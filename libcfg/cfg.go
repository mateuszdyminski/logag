package libcfg

import (
	"flag"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

var (
	configPath        = flag.String("config", "config/logag_config.toml", "path to logag configuration")
	staticsPath       = flag.String("statics", "", "Path to directory with statics")
	host              = flag.String("host", "", "host address")
	httpDrainInterval = flag.String("http-drain-interval", "", "Http drain interval")
	batchSize         = flag.Int("batch-size", 3, "ElasticSearch ingres batch size")
	boltDDPath         = flag.String("bolt-db-path", "", "Path to boltDB data file")
)

type Cfg struct {
	Elastics          []string
	Host              string
	HttpDrainInterval string
	StaticsPath       string
	BatchSize         int
	BoltDDPath        string
}

func LoadCfg() (*Cfg, error) {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*configPath)
	if err != nil {
		return nil, err
	}

	conf := Cfg{}
	if err := toml.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}

	if *staticsPath != "" {
		conf.StaticsPath = *staticsPath
	}

	if *host != "" {
		conf.Host = *host
	}

	if *httpDrainInterval != "" {
		conf.HttpDrainInterval = *httpDrainInterval
	}

	if *boltDDPath != "" {
		conf.BoltDDPath = *boltDDPath
	}

	conf.BatchSize = *batchSize

	return &conf, nil
}
