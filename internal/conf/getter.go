package conf

import "flag"

var (
	exchange string
	endpoint string
	logLevel string
	confFile string
)

type Conf struct {
	Endpoint string
	Exchange []string
	LogLevel string
}

func init() {
	flag.StringVar(&endpoint, "endpoint", "", "endpoint ws address")
	flag.StringVar(&endpoint, "exchange", "", "exchange")
	flag.StringVar(&logLevel, "log_level", "info", "log level")
	flag.StringVar(&confFile, "conf", "", "path of the configuration file")
}

func Get() Conf {
	flag.Parse()

	return Conf{
		Endpoint: endpoint,
	}
}
