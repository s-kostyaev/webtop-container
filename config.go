package main

import (
	"github.com/BurntSushi/toml"
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"time"
)

const (
	configPath        = "/etc/webtop/container.toml"
	comment           = "webtop"
	staticDir         = "/usr/share/webtop/"
	dataSocketPath    = "/webtop-data.sock"
	commandSocketPath = "/webtop-command.sock"
)

var (
	logfile   = os.Stderr
	formatter = logging.MustStringFormatter(
		"%{time:15:04:05.000000} %{pid} %{level:.8s} %{longfile} %{message}")
	loglevel = logging.INFO
	logger   = logging.MustGetLogger("webtop")
)

func setupLogger() {
	logging.SetBackend(logging.NewLogBackend(logfile, "", 0))
	logging.SetFormatter(formatter)
	logging.SetLevel(loglevel, logger.Module)
}

type Config struct {
	WaitTimeout duration
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func getConfig(configPath string) *Config {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Fatal(err.Error())
	}
	config := Config{}
	_, err = toml.Decode(string(buf), &config)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return &config
}
