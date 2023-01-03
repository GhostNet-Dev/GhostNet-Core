package gconfig

const (
	DefaultBinaryVersion = 1
	DefaultIp            = "127.0.0.1"
	DefaultPort          = "50129"
	DefaultRootPath      = "./"
	DefaultSqlPath       = "./"
)

var (
	GlobalConfig GConfig
)

type GConfig struct {
	GhostVersion uint32
	Ip           string
	Port         string
	RootPath     string
	SqlPath      string
}

func DefaultConfig() *GConfig {
	GlobalConfig = GConfig{
		GhostVersion: 1,
		Ip:           DefaultIp,
		Port:         DefaultPort,
		RootPath:     DefaultRootPath,
		SqlPath:      DefaultSqlPath,
	}
	return &GlobalConfig
}
