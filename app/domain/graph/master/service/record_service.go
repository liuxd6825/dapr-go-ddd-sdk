package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
)

const (
	defaultMinioName    = "default"
	defaultBucketConfig = "importData"
	defaultBatchSize    = 500
	defaultRetries      = 1
	defaultConcurrency  = 1
)

// RecordParquetService
// @Description: 银行流水 parquet 导入业务编排
type RecordParquetService struct {
	dao *dao.RecordParquetDao
	env *env.Env
}

func NewRecordParquetService(env *env.Env) *RecordParquetService {
	return &RecordParquetService{
		dao: dao.NewRecordParquetDao(),
		env: env,
	}
}

// ImportOptions 业务层可选参数
type ImportOptions struct {
	BatchSize   int
	Parallel    bool
	Concurrency int
	Retries     int
}

// ImportFromS3 导入指定 s3Path 到 Neo4j (按日期聚合)
// s3Path 支持三种形式:
//  1. "s3://user:pass@host/bucket/key"  完整 URL, 直接使用
//  2. "s3://bucket/key"                 简写, 从 env minio 注入凭证
//  3. "bucket/key"                      简写, 走 env minio + default 桶
//
// caseId 必传: 用于构建 label 与关系属性
func (s *RecordParquetService) ImportFromS3(ctx context.Context, s3Path, caseId string, opts ImportOptions) (*dao.ImportResult, error) {
	if s3Path == "" {
		return nil, errors.New("s3Path is empty")
	}
	if caseId == "" {
		return nil, errors.New("caseId is empty")
	}

	tenantId := appctx.GetTenantId2(ctx)
	if tenantId == "" {
		return nil, errors.New("tenantId is empty in context")
	}

	fullURL, err := s.resolveS3URL(s3Path)
	if err != nil {
		return nil, err
	}

	params := dao.ImportParams{
		TenantId:    tenantId,
		CaseId:      caseId,
		S3Path:      fullURL,
		BatchSize:   opts.BatchSize,
		Parallel:    opts.Parallel,
		Concurrency: opts.Concurrency,
		Retries:     opts.Retries,
	}
	if params.BatchSize <= 0 {
		params.BatchSize = defaultBatchSize
	}
	if params.Retries <= 0 {
		params.Retries = defaultRetries
	}
	if params.Concurrency <= 0 {
		params.Concurrency = defaultConcurrency
	}

	var res *dao.ImportResult
	gp.Try(func() error {
		r, e := s.dao.ImportFromS3(ctx, params)
		if e != nil {
			return e
		}
		res = r
		return nil
	}).Catch(func(e error) {
		err = e
	}).Finally(func() {})

	return res, err
}

// resolveS3URL 把简写路径转成 s3://user:pass@host:port/bucket/key
func (s *RecordParquetService) resolveS3URL(s3Path string) (string, error) {
	// 形式 1: s3://user:pass@host/bucket/key  已是完整 URL
	if strings.HasPrefix(s3Path, "s3://") {
		trimmed := strings.TrimPrefix(s3Path, "s3://")
		// 已包含 @  (user:pass@) -> 直接使用
		if strings.Contains(trimmed, "@") {
			return s3Path, nil
		}
		// 形式 2: s3://host/bucket/key -> 注入 user:pass@
		minio, ok := s.env.GetMinioByKey(defaultMinioName)
		if !ok {
			return "", errors.New("env minio.default not configured")
		}
		parts := strings.SplitN(trimmed, "/", 3)
		if len(parts) < 3 {
			return "", errors.New("s3 path must include host/bucket/key, got: %s", s3Path)
		}
		return fmt.Sprintf("s3://%s:%s@%s/%s/%s",
			minio.AccessKey, minio.SecretKey, parts[0], parts[1], parts[2]), nil
	}

	// 形式 3: bucket/key -> 走 env minio + default 桶
	minio, ok := s.env.GetMinioByKey(defaultMinioName)
	if !ok {
		return "", errors.New("env minio.default not configured")
	}
	bucketName, ok := minio.Buckets[defaultBucketConfig]
	if !ok {
		return "", errors.New("env minio.default.buckets.importData not configured")
	}
	return dao.BuildS3URL(minio, bucketName, s3Path), nil
}
