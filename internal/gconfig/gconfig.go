package gconfig

const (
	DefaultBinaryVersion  = 1
	DefaultConfigFilename = "global"
	EnvPrefix             = "GON" // GhOstNet

	DefaultIp                  = "127.0.0.1"
	DefaultPort                = "50129"
	DefaultGrpcPort            = "50229"
	DefaultRootPath            = "./"
	DefaultSqlPath             = "./"
	DefaultFilePath            = "./"
	DefaultDbName              = "sqlite3"
	DefaultDbLiteStoreFilename = "litestore.db"

	DefaultStandaloneMode = false
	DefaultTimeout        = 8
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
	LiteStoreFilename     string
	FilePath              string
	DbName                string
	StandaloneMode        bool
	Id                    uint32

	Timeout           uint32
	BlockTickInterval int
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
		LiteStoreFilename:     DefaultDbLiteStoreFilename,
		FilePath:              DefaultFilePath,
		DbName:                DefaultDbName,
		StandaloneMode:        DefaultStandaloneMode,
		Timeout:               DefaultTimeout,
		BlockTickInterval:     DefaultTimeout,
	}
}
