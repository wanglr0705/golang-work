package config

// 配置信息
type Config struct {
	Host     string         `yaml:"Host"`
	Port     string         `yaml:"Port"`
	Database DatabaseConfig `yaml:"Database"`
	Redis    RedisConfig    `yaml:"Redis"`
}

// mysql数据库的配置信息
type DatabaseConfig struct {
	UserName  string `yaml:"UserName"`
	PassWord  string `yaml:"PassWord"`
	IdAddress string `yaml:"IdAddress"`
	Port      string `yaml:"Port"`
	DbName    string `yaml:"DbName"`
	Charset   string `yaml:"Charset"`
}

// Redis的配置信息
type RedisConfig struct {
	Network  string
	Address  string
	Password string
}
