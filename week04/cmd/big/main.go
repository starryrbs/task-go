package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/task-go/week04/internal/conf"
	"gopkg.in/yaml.v3"
	"os"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "configs/config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	bytes, err := os.ReadFile(flagconf)
	if err != nil {
		panic(err)
	}

	var c conf.Conf
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", c)

	app, err := initApp(context.Background(), &c)
	if err != nil {
		panic(err)
	}
	app.Run(c.Server.Address)
}
