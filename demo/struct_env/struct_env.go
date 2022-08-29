package main

import "github.com/delfanhao/go-conf/conf"

type StructConfig struct {
	FromJson struct {
		Item1   string
		Item2   string
		Port    int
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

	FromEnv []int
}

func main() {
	cfg := StructConfig{}

	conf.TRACE = true
	conf.Load(&cfg)

	println("FromJson.Item1 = ", cfg.FromJson.Item1)
	println("FromJson.Prot = ", cfg.FromJson.Port)
	println("FromJson.SubItem.Item1 = ", cfg.FromJson.SubItem.Item1)
	println("FromYml.Item1 = ", cfg.FromYml.Item1)
	println("FromYml.SubItem.Item1 = ", cfg.FromYml.SubItem.Item1)

	if len(cfg.FromEnv) == 0 {
		println("FromEnv is not found.")
	}
	for s := range cfg.FromEnv {
		println("FromEnv = ", s, cfg.FromEnv[s])
	}

}
