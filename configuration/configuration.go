package configuration

import (
	"fmt"
	"io/ioutil"
	"time"

	"os"

	"github.com/marguerite/util/dir"
	"github.com/marguerite/util/slice"
	//yaml "gopkg.in/yaml.v2"
	simplejson "github.com/bitly/go-simplejson"
)

// Configurations collection of individual configurations
type Configurations []*Configuration

// Configuration basic configuration structure
/*type Configuration struct {
	Name            string `yaml:"name"`
	Version         string `yaml:"version"`
	PreviousVersion string `yaml:"previous"`
	Plugin          string `yaml:"plugin"`
	URL             string `yaml:"url"`
	Maintainer      string `yaml:"maintainer"`
	Email           string `yaml:"email"`
	Unstable        bool   `yaml:"unstable"`
	GenChange       bool   `yaml:"genchange"`
	MTime           string `yaml:"mtime"`
	NonDigit        bool   `yaml:"nondigit"`
}*/
type Configuration simplejson.Json

func (c *Configuration) StringAttr(s string) string {
	json, ok := (*simplejson.Json)(c).CheckGet(s)
	if !ok {
		fmt.Printf(".json doesn't contain attribute %s\n", s)
		os.Exit(1)
	}
	return json.MustString()
}

func (c *Configuration) BoolAttr(s string) bool {
	json, ok := (*simplejson.Json)(c).CheckGet(s)
	if !ok {
		fmt.Printf(".json doesn't contain attribute %s\n", s)
		os.Exit(1)
	}
	return json.MustBool()
}

// ModificationTime modification time in golang time
func (c *Configuration) ModificationTime() time.Time {
	timeForm := "2006-01-02 15:04:05"
	t, err := time.Parse(timeForm, c.StringAttr("mtime"))
	if err != nil {
		fmt.Println("time format in .json is wrong, correct format: 2020:03:31 13:27")
		os.Exit(1)
	}
	return t
}

// ParseJson parse .json file in config directory
func ParseJson() Configurations {
	config := Configurations{}
	jsons, err := dir.Ls("./config/*.json")

	if err != nil {
		fmt.Println("Can not find any .json configuration in config directory.")
		os.Exit(1)
	}

	for _, v := range jsons {
		b, err := ioutil.ReadFile(v)
		if err != nil {
			fmt.Printf("Can not read file %s\n", v)
			os.Exit(1)
		}
		json, err := simplejson.NewJson(b)
		if err != nil {
			fmt.Printf("Can not load json configuration %s\n", v)
			os.Exit(1)
		}
		slice.Concat(&config, (*Configuration)(json))
	}
	return config
}
