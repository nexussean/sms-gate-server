package db

var ConfigDefault = Config{
	Host:     "localhost",
	Port:     3306,
	User:     "root",
	Password: "",
	Database: "db",
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// Helper function to set default values
func configDefault(config Config) Config {
	// Override default config
	if config.Host == "" {
		config.Host = ConfigDefault.Host
	}
	if config.Port == 0 {
		config.Port = ConfigDefault.Port
	}
	if config.User == "" {
		config.User = ConfigDefault.User
	}
	if config.Password == "" {
		config.Password = ConfigDefault.Password
	}
	if config.Database == "" {
		config.Database = ConfigDefault.Database
	}

	return config
}
