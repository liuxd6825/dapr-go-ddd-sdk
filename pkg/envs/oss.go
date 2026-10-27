package envs

import (
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/boolutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils"
)

type OssType string

const (
	OssTypeMinio  OssType = "minio"
	OssTypeS3     OssType = "s3"
	OssTypeRustFS OssType = "rustfs"
)

type OssEnv struct {
	items map[string]OSSItem `yaml:"-"`
	minio map[string]*Minio  `yaml:"-"`
}

func newOssEnv() *OssEnv {
	return &OssEnv{
		items: map[string]OSSItem{},
		minio: map[string]*Minio{},
	}
}

func initOSS(env *Env) {
	errList := errors.NewErrors()
	const ossTypeConst = "ossType"
	const enableConst = "enable"
	for fsKey, info := range env.OSS {
		fsType, err := maputils.GetString(info, ossTypeConst, "")
		if err != nil {
			errList.AddError(errors.New("%s fsType is required", err.Error()))
		}
		fsType = strings.ToLower(fsType)
		info[ossTypeConst] = fsType
		enable := boolutils.GetBool(info[enableConst], true)
		if !enable {
			continue
		}
		switch fsType {
		case OssTypeMinio.String():
			oss, err := newMinioWithMap(info)
			if err != nil {
				errList.AddError(err)
				continue
			}
			if err := addOssMinio(env, oss); err != nil {
				errList.AddError(err)
			}
			env.oss.minio[oss.OssKey] = oss
			break
		case OssTypeRustFS.String():
			errList.AddError(errors.New("%s ossType not supported %s", fsType, OssTypeRustFS.String()))
			break
		case OssTypeS3.String():
			errList.AddError(errors.New("%s ossType not supported %s", fsType, OssTypeS3.String()))
			break
		default:
			errList.AddError(errors.New("%s ossType is required", fsKey))
		}
	}
}

func (d OssType) String() string {
	return string(d)
}
