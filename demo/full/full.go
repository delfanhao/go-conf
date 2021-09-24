package main

import "github.com/delfanhao/go-conf/conf"

type FullConfig struct {
	FromApplication string `default:"Value from tag default"`
	FromTagDefault  string `default:"Value from tag default"`

	FromConfAppYml      string
	FromConfAppJson     string
	FromConfAppIni      string
	FromAppYml          string
	FromAppJson         string
	FromAppIni          string
	FromConfDefaultYml  string
	FromConfDefaultJson string
	FromConfDefaultIni  string
	FromDefaultYml      string
	FromDefaultJson     string
	FromDefaultIni      string
	FromCmdLine         string
	FromEnv             string

	GeneralItem string
}

func main() {
	conf.TRACE = true
	cfg := FullConfig{
		FromApplication: "Value from application",
	}

	conf.Load(&cfg)

	println("FromApplication value = ", cfg.FromApplication)
	println("FromTagDefault value = ", cfg.FromTagDefault)

	println("FromConfAppYml value = ", cfg.FromConfAppYml)
	println("FromConfAppJson value = ", cfg.FromConfAppJson)
	println("FromConfAppIni value = ", cfg.FromConfAppIni)
	println("FromAppYml value = ", cfg.FromAppYml)
	println("FromAppJson value = ", cfg.FromAppJson)
	println("FromAppIni value = ", cfg.FromAppIni)
	println("FromConfDefaultYml value = ", cfg.FromConfDefaultYml)
	println("FromConfDefaultJson value = ", cfg.FromConfDefaultJson)
	println("FromConfDefaultIni value = ", cfg.FromConfDefaultIni)
	println("FromDefaultYml value = ", cfg.FromDefaultYml)
	println("FromDefaultJson value = ", cfg.FromDefaultJson)
	println("FromDefaultIni value = ", cfg.FromDefaultIni)
	println("FromCmdLine value = ", cfg.FromCmdLine)
	println("FromEnv value = ", cfg.FromEnv)

	println("GeneralItem value = ", cfg.GeneralItem)

}
