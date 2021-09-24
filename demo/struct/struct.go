package main

import "github.com/delfanhao/go-conf/conf"

type StructConfig struct {
	FromJson struct {
		Item1   string
		Item2   string
		SubItem struct {
			Item1 string
			Item2 string
		}
	}

	FromYml struct {
		Item1   string
		Item2   string
		SubItem struct {
			Item1 string
			Item2 string
		}
	}
}

func main() {
	cfg := StructConfig{}

	conf.TRACE = true
	conf.Load(&cfg)

	println("FromJson.Item1 = ", cfg.FromJson.Item1)
	println("FromJson.SubItem.Item1 = ", cfg.FromJson.SubItem.Item1)
	println("FromYml.Item1 = ", cfg.FromYml.Item1)
	println("FromYml.SubItem.Item1 = ", cfg.FromYml.SubItem.Item1)
}
