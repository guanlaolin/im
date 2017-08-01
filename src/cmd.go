package main

import (
	"flag"
)

var help = flag.String("help", "", `
		-help		帮助
		-config		设置配置文件路径
	`)

var config = flag.String("config", "../config/config.json", "-config='config/config.json'")
