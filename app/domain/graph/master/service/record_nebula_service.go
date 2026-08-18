package service

import (
	"context"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
)

const (
	nebulaMinioName    = "default"
	nebulaBucketConfig = "importData"
	nebulaBatchSize    = 500
	nebulaRetries      = 1
	nebulaConcurrency  = 1
)

// RecordNebulaService 银行流水 parquet -> NebulaGraph 导入业务编排
type RecordNebulaService struct {
	dao *dao.RecordNebulaDao
	env *env.Env
}

func NewRecordNebulaService(env *env.Env) *RecordNebulaService {
	return &RecordNebulaService{
		dao: dao.NewRecordNebulaDao(env),
		env: env,
	}
}

// ImportFromS3 导入指定 s3Path 到 NebulaGraph
func (s *RecordNebulaService) ImportFromS3(ctx context.Context, s3Path, caseId string, opts ImportOptions) (*dao.ImportResult, error) {
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

	bucket, key, err := s.resolveS3Path(s3Path)
	if err != nil {
		return nil, err
	}

	params := dao.NebulaImportParams{
		TenantId:    tenantId,
		CaseId:      caseId,
		Bucket:      bucket,
		Key:         key,
		BatchSize:   opts.BatchSize,
		Parallel:    opts.Parallel,
		Concurrency: opts.Concurrency,
		Retries:     opts.Retries,
	}
	if params.BatchSize <= 0 {
		params.BatchSize = nebulaBatchSize
	}
	if params.Retries <= 0 {
		params.Retries = nebulaRetries
	}
	if params.Concurrency <= 0 {
		params.Concurrency = nebulaConcurrency
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

// resolveS3Path 解析三种形式返回 (bucket, key)
func (s *RecordNebulaService) resolveS3Path(s3Path string) (bucket, key string, err error) {
	if s3Path == "" {
		return "", "", errors.New("s3Path is empty")
	}
	if strings.HasPrefix(s3Path, "s3://") {
		trimmed := strings.TrimPrefix(s3Path, "s3://")
		if strings.Contains(trimmed, "@") {
			// 形式 1: s3://user:pass@host/bucket/key
			at := strings.Index(trimmed, "@")
			rest := trimmed[at+1:]
			parts := strings.SplitN(rest, "/", 3)
			if len(parts) != 3 {
				return "", "", errors.New("s3 path must include host/bucket/key, got: %s", s3Path)
			}
			return parts[1], parts[2], nil
		}
		// 形式 2: s3://bucket/key
		minioCfg, ok := s.env.GetMinioByKey(nebulaMinioName)
		if !ok {
			return "", "", errors.New("env minio.default not configured")
		}
		_ = minioCfg
		parts := strings.SplitN(trimmed, "/", 2)
		if len(parts) != 2 {
			return "", "", errors.New("s3 path must include bucket/key, got: %s", s3Path)
		}
		return parts[0], parts[1], nil
	}

	// 形式 3: bucket/key
	minioCfg, ok := s.env.GetMinioByKey(nebulaMinioName)
	if !ok {
		return "", "", errors.New("env minio.default not configured")
	}
	bucketName, ok := minioCfg.Buckets[nebulaBucketConfig]
	if !ok {
		return "", "", errors.New("env minio.default.buckets.importData not configured")
	}
	return bucketName, s3Path, nil
}