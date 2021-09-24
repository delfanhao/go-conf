package main

import (
	"github.com/delfanhao/go-conf/conf"
)

// DemoConfig 定义一个Demo的配置
type DemoConfig struct {
	GeneralStringItem string
	GeneralIntItem    int
	GeneralFloatItem  float64
}

func main() {
	conf.TRACE = false
	cfg := DemoConfig{}
	conf.Load(&cfg)
}
