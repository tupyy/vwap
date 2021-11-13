// This package parse the configuration flags or the conf json file.
package conf

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/tupyy/vwap/internal/log"
)

var (
	pairs         string
	endpoint      string
	logLevel      string
	confFile      string
	outputFile    string
	maxDataPoints int64
)

type Conf struct {
	Endpoint      string
	TradingPairs  []string
	MaxDataPoints int64
	OutputFile    string
}

func init() {
	flag.StringVar(&endpoint, "endpoint", "", "endpoint ws address")
	flag.StringVar(&pairs, "pairs", "", "comma separated trading pairs")
	flag.StringVar(&logLevel, "log_level", "info", "log level")
	flag.StringVar(&confFile, "config", "", "path of the configuration file")
	flag.StringVar(&outputFile, "output", "", "path of the output file")
	flag.Int64Var(&maxDataPoints, "max_data_points", 200, "maximum number of data points used to compute the average")
}

func Get() Conf {
	flag.Parse()

	// if conf is set, ignore the rest of the flags and parse the configuration file
	if len(confFile) > 0 {
		content, err := os.ReadFile(confFile)
		if err != nil {
			panic(err)
		}

		conf := parseConfFile(content)

		if conf.MaxDataPoints == 0 {
			log.GetLogger().Warning("cannot set max data points to 0. Default to 200.")

			conf.MaxDataPoints = 200
		}

		return conf
	}

	// check if endpoint and trading pairs are set
	if len(endpoint) == 0 || len(pairs) == 0 {
		log.GetLogger().Error("Both endpoint and trading pairs are mandatory.")

		os.Exit(1)
	}

	conf := Conf{
		Endpoint:     endpoint,
		TradingPairs: make([]string, 0, 3),
		OutputFile:   outputFile,
	}

	for _, p := range strings.Split(pairs, ",") {
		conf.TradingPairs = append(conf.TradingPairs, p)
	}

	if maxDataPoints == 0 {
		log.GetLogger().Warning("cannot set max data points to 0. Default to 200.")

		maxDataPoints = 200
	}

	conf.MaxDataPoints = maxDataPoints

	log.SetLogLevel(parseLogLevel(logLevel))

	return conf
}

func parseLogLevel(l string) log.Level {
	switch strings.ToLower(l) {
	case "trace":
		return log.Trace
	case "debug":
		return log.Debug
	case "warning":
		return log.Warning
	case "error":
		return log.Error
	default:
		return log.Info
	}
}

func parseConfFile(content []byte) Conf {
	confFile := struct {
		Endpoint      string   `json:"endpoint"`
		TradingPairs  []string `json:"trading_pairs"`
		LogLevel      string   `json:"log_level,omitempty"`
		MaxDataPoints int64    `json:"max_data_points,omitempty"`
		OutputFile    string   `json:"output_file,omitempty"`
	}{}

	// unmarshal the content into confFile
	if err := json.Unmarshal(content, &confFile); err != nil {
		panic(err)
	}

	log.SetLogLevel(parseLogLevel(confFile.LogLevel))

	return Conf{
		Endpoint:      confFile.Endpoint,
		TradingPairs:  confFile.TradingPairs,
		MaxDataPoints: confFile.MaxDataPoints,
		OutputFile:    confFile.OutputFile,
	}
}
