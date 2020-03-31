package main

import (
	"fmt"
	"path/filepath"
	"plugin"

	"os"

	"github.com/marguerite/util/dir"
	"github.com/marguerite/util/slice"
	"github.com/openSUSE/sauria/commit"
	"github.com/openSUSE/sauria/configuration"
)

func main() {
	cwd, _ := os.Getwd()
	configs := configuration.ParseYAML()
	plugins, err := dir.Ls("./plugins/*.so")
	if err != nil || len(plugins) == 0 {
		fmt.Println("Can not find any plugin in plugins directory. Forgot to build them first?")
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
			fn := fetch.(func(configuration.Configuration) (commit.Commit, error))
			release, err := fn(c)
			if err != nil {
				panic(err)
			}
			fmt.Println(release.UnstableVersion())
		}
	}
}
