package xtest

type TestType string

const (
	TestType_Mongo224   = "Mongo224"
	TestType_MongoLocal = "MongoLocal"
)

const dbKey = "$db"

func Init(testType TestType) {
	InitTimeZone()
	switch testType {
	case TestType_Mongo224:
		opts := &MongoOptions{
			DBKey:      str(MongoDBKey),
			Host:       str(MongoHostRemote),
			DBName:     str("master"),
			ReplicaSet: str("mongors"),
			User:       str("super_admin"),
			Pwd:        str("123456"),
		}
		InitEnv_MongoRemoteTest(opts)
	case TestType_MongoLocal:
		opts := &MongoOptions{
			DBKey:      str(MongoDBKey),
			Host:       str(MongoHostLocal),
			DBName:     str("master"),
			ReplicaSet: str("rs0"),
			User:       str("super_admin"),
			Pwd:        str("123456"),
		}
		InitEnv_MongoRemoteTest(opts)
	}
}

func str(str string) *string {
	return &str
}

/*
// Init
// @Description: 初始化QueryService测试环境
func Init(ctx context.Context, fileName string, envName string, eventTypes []restapp2.RegisterEventType) error {
	if fileName == "" {
		return errors.New("fileName is null")
	}
	if envName == "" {
		return errors.New("envName is null")
	}
	// 加载配置文件
	config, err := restapp2.NewConfigByFile("${search}/config/" + fileName)
	if err != nil {
		return err
	}

	// 读取指定环境下的配置信息
	env, err := config.GetEnvConfig(envName)
	if err != nil {
		return err
	}

	return restapp2.InitApplication(ctx, env, eventTypes, true, nil)
}
*/
