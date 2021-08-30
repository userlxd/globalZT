package config

var Config config

type config struct {
	Mode    string
	LogFile string
}

func init() {
	Config = config{Mode: "dev"}
	Config.loadFlag()
}
