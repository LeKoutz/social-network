package utils

var GlobalConfig Configuration = DefaultConfiguration()

type Configuration struct {
	DbPath string
	Ip     string
	Port   string
	Debug  bool
}

func DefaultConfiguration() Configuration {
	return Configuration{
		DbPath: "./db.db",
		Ip:     "127.0.0.1",
		Port:   "8080",
		Debug:  false,
	}
}
