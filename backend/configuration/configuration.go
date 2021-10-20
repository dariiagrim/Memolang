package configuration

var configuration *Configuration

type Configuration struct {
	ServerConfig ServerConfig
	StoreConfig  StoreConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type StoreConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

func GetConfiguration() *Configuration {
	if configuration != nil {
		return configuration
	}
	configuration = &Configuration{
		ServerConfig: ServerConfig{
			Host: "localhost",
			Port: ":8080",
		},
		StoreConfig: StoreConfig{
			Host:     "localhost",
			User:     "postgres",
			Password: "qwerty",
			Name:     "memolang",
			Port:     "5432",
		},
	}

	return configuration
}
