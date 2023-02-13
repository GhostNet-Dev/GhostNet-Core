package gconfig

const (
	DefaultBinaryVersion  = 1
	DefaultConfigFilename = "global"
	EnvPrefix             = "GON" // GhOstNet

	DefaultIp           = "127.0.0.1"
	DefaultPort         = "50129"
	DefaultGrpcPort     = "50229"
	DefaultRootPath     = "./"
	DefaultSqlPath      = "./"
	DefaultFilePath     = "./"
	DefaultDbName       = "sqlite3"
	DefaultDbScheme     = "db.sqlite3.sql"
	DefaultDbSchemePath = "./"

	DefaultStandaloneMode = false
)

type GConfig struct {
	GhostVersion uint32
	Username     string
	Password     []byte

	DefaultConfigFilename string
	EnvPrefix             string
	Ip                    string
	Port                  string
	GrpcPort              string
	RootPath              string
	SqlPath               string
	FilePath              string
	DbName                string
	DbScheme              string
	DbSchemePath          string
	StandaloneMode        bool

	Timeout uint32
}

func NewDefaultConfig() *GConfig {
	return &GConfig{
		GhostVersion:          1,
		DefaultConfigFilename: DefaultConfigFilename,
		EnvPrefix:             EnvPrefix,
		Ip:                    DefaultIp,
		Port:                  DefaultPort,
		GrpcPort:              DefaultGrpcPort,
		RootPath:              DefaultRootPath,
		SqlPath:               DefaultSqlPath,
		FilePath:              DefaultFilePath,
		DbName:                DefaultDbName,
		DbScheme:              DefaultDbScheme,
		DbSchemePath:          DefaultDbSchemePath,
		StandaloneMode:        DefaultStandaloneMode,
	}
}
