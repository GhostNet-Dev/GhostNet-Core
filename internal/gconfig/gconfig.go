package gconfig

const (
	DefaultBinaryVersion  = 1
	DefaultConfigFilename = "global"
	EnvPrefix             = "GON" // GhOstNet

	DefaultIp           = "127.0.0.1"
	DefaultPort         = "50129"
	DefaultGrpcPort     = "50130"
	DefaultRootPath     = "./"
	DefaultSqlPath      = "./"
	DefaultFilePath     = "./"
	DefaultDbName       = "sqlite3"
	DefaultDbScheme     = "db.sqlite3.sql"
	DefaultDbSchemePath = "./"

	DefaultStandaloneMode = false
)

var (
	GlobalConfig *GConfig
)

type GConfig struct {
	GhostVersion uint32
	Username     string
	Password     string

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
}

func DefaultConfig() *GConfig {
	if GlobalConfig == nil {
		GlobalConfig = &GConfig{
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
	return GlobalConfig
}
