package main

import "github.com/delfanhao/go-conf/conf"

type SliceConfig struct {
	Ips []string
}

func main() {
	cfg := SliceConfig{}
	conf.TRACE = true
	conf.Load(&cfg)
	println(len(cfg.Ips))
	for i := range cfg.Ips {
		println(cfg.Ips[i])
	}
}
