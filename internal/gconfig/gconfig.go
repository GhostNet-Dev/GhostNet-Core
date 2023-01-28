package gconfig

const (
	DefaultBinaryVersion = 1
	DefaultIp            = "127.0.0.1"
	DefaultPort          = "50129"
	DefaultRootPath      = "./"
	DefaultSqlPath       = "./"
	DefaultFilePath      = "./"
	DefaultDbName        = "sqlite3"
	DefaultDbScheme      = "db.sqlite3.sql"
	DefaultDbSchemePath  = "./"
)

var (
	GlobalConfig GConfig
)

type GConfig struct {
	GhostVersion   uint32
	Ip             string
	Port           string
	RootPath       string
	SqlPath        string
	FilePath       string
	DbPath         string
	DbName         string
	StandaloneMode bool
}

func DefaultConfig() *GConfig {
	GlobalConfig = GConfig{
		GhostVersion: 1,
		Ip:           DefaultIp,
		Port:         DefaultPort,
		RootPath:     DefaultRootPath,
		SqlPath:      DefaultSqlPath,
		FilePath:     DefaultFilePath,
		DbName:       DefaultDbName,
	}
	return &GlobalConfig
}
