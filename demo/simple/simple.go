package main

import (
	"fmt"
	"github.com/delfanhao/go-conf/conf"
)

type SimpleStruct struct {
	Host string
	Port int
}

func main() {
	cfg := SimpleStruct{
		Port: 80,
	}
	conf.TRACE = true
	conf.Load(&cfg)

	println(fmt.Sprintf("Host:%s,port:%d", cfg.Host, cfg.Port))
}
