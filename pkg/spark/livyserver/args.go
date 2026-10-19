package livyserver

type SparkOption struct {
	DriverMemory             string // "8g",
	ExecutorMemory           string // "8g",
	ExecutorCores            int    //  6,
	JavaJarFile              string // Java Jar包文件
	JavaApplicationClassName string // Spark加载类名
	S3                       S3Option
}

type S3Option struct {
	Endpoint string // "http://192.168.120.224:9000"
	User     string //
	Password string //
}

type MongoOption struct {
}

type S3JarOption struct {
	JavaJarFile   string
	JavaClassName string
	Endpoint      string // "http://192.168.120.224:9000"
	User          string //
	Password      string //
}

type MySQLOption struct {
	Host     string
	Port     string
	DB       string
	Table    string
	User     string
	Password string
}
