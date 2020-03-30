package main

import (
	"fmt"
	"io/ioutil"

	"os"

	"github.com/marguerite/util/dir"
	"github.com/marguerite/util/slice"
	yaml "gopkg.in/yaml.v2"
)

// Configurations collection of individual configurations
type Configurations []Configuration

// Configuration basic configuration structure
type Configuration struct {
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	Plugin     string `yaml:"plugin"`
	URL        string `yaml:"url"`
	Maintainer string `yaml:"maintainer"`
	Email      string `yaml:"email"`
	Unstable   bool   `yaml:"unstable"`
	GenChange  bool   `yaml:"genchange"`
}

func parseYAML() Configurations {
	config := Configurations{}
	ymls, err := dir.Ls("./config/*.yml")

	if err != nil {
		fmt.Println("Can not find any .yml configuration in config directory.")
		os.Exit(1)
	}

	for _, v := range ymls {
		c := Configurations{}
		b, err := ioutil.ReadFile(v)
		if err != nil {
			fmt.Printf("Can not read file %s\n", v)
			os.Exit(1)
		}
		err = yaml.Unmarshal(b, &c)
		if err != nil {
			fmt.Printf("Can not load yaml configuration %s\n", v)
			os.Exit(1)
		}
		slice.Concat(&config, c)
	}
	return config
}
