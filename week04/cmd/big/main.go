package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "configs/config.yaml", "config path, eg: -conf config.yaml")
}

type Config struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

func main() {
	flag.Parse()

	bytes, err := os.ReadFile(flagconf)
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", cfg)
}
