package main

import (
	conf2 "go-conf/conf"
)

type DevConfig struct {
	Keys []int
	Key2 []string
}

func main() {
	conf2.TRACE = true
	cfg := &DevConfig{}
	conf2.Load(cfg)
	println(len(cfg.Keys), len(cfg.Key2))
	for i := range cfg.Keys {
		println(cfg.Keys[i])
	}
	for i := range cfg.Key2 {
		println(cfg.Key2[i])
	}
}
