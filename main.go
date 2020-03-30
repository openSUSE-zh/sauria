package main

import (
	"fmt"
	"path/filepath"
	"plugin"

	"os"

	"github.com/marguerite/util/dir"
	"github.com/marguerite/util/slice"
	"github.com/openSUSE/sauria/common"
)

func main() {
	cwd, _ := os.Getwd()
	configs := parseYAML()
	plugins, err := dir.Ls("./plugins/*.so")
	if err != nil {
		fmt.Println("Can not find any plugin in plugins directory.")
		os.Exit(1)
	}

	for _, c := range configs {
		p := filepath.Join(cwd, "plugins", c.Plugin+".so")
		if ok, err := slice.Contains(plugins, p); ok && err == nil {
			plug, err := plugin.Open(p)
			if err != nil {
				fmt.Printf("Can not open plugin %s\n", p)
				os.Exit(1)
			}
			fetch, err := plug.Lookup("FetchNewVersion")
			if err != nil {
				fmt.Printf("Can not find function FetchNewVersion in plugin %s\n", p)
				os.Exit(1)
			}
			fn := fetch.(func(string, bool, bool) (common.Commit, error))
			fmt.Println(fn(c.URL, c.Unstable, c.GenChange))
		}
	}
}
