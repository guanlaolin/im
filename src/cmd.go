package main

import (
	"flag"
)

var help = flag.String("help", "", `
		--help		show help
		--config	set config file path
	`)

var config = flag.String("config", "../config/config.json", "set config file path")
