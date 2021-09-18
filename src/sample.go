package main

import (
	"go-conf/src/conf"
)

type Config struct {
	Ftp struct {
		Host string `default:"0.0.0.0"`
		Port int
	}

	Root string

	Target struct {
		File struct {
			Name string
			Size int
		}

		Storage string `default:"Default of target.storage"`
	}
}

func main() {
	cfg := Config{
		Root: "aaa",
	}
	conf.Load(&cfg)
	show(&cfg)
}

func show(cfg *Config) {
	println("- START -")
	println("ftp.host = ", cfg.Ftp.Host)
	println("ftp.port = ", cfg.Ftp.Port)
	println("root = ", cfg.Root)
	println("target.file.name = ", cfg.Target.File.Name)
	println("target.file.size = ", cfg.Target.File.Size)
	println("target.storage = ", cfg.Target.Storage)
	println("- END- ")
}
