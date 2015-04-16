package main

import (
	"github.com/op/go-logging"
	"os"
)

const (
	configPath = "/etc/webtop/container.toml"
	comment    = "webtop"
	staticDir  = "/usr/share/webtop/"
	socketPath = "/webtop.sock"
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
