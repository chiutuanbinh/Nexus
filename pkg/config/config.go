package config

type Config struct {
	ServerPort    string `mapstructure:"server_port"`
	StorageDir    string `mapstructure:"storage_dir"`
	SegmentSizeMB int    `mapstructure:"segment_size_mb"`
	MemtableMaxMB int    `mapstructure:"memtable_max_mb"`
	UseHash       bool   `mapstructure:"use_hash"`
}

var configInstance *Config

func Init(initFunc func(*Config) error) {
	var config Config
	err := initFunc(&config)
	if err != nil {
		panic(err)
	}
	configInstance = &config
}

func GetConfig() *Config {
	return configInstance
}
