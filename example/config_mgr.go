package example

import "fmt"

var ConfigMgr = &configMgr{
	Configs: map[string]Config{},
}

type (
	Config interface {
		Load(filepath string)
		Filename() string
	}

	configMgr struct {
		Configs map[string]Config
	}
)

func (c *configMgr) Init() {
	c.Register(TestSheet1)

}

func (c *configMgr) Register(conf Config) {
	if _, ok := c.Configs[conf.Filename()]; ok {
		panic("error: duplicate config " + conf.Filename())
	}

	c.Configs[conf.Filename()] = conf
}

func (c *configMgr) Load(configPath string) {
	fmt.Println("Load Configs ......")
	for _, conf := range c.Configs {
		conf.Load(configPath)
		fmt.Println("load config " + conf.Filename() + " SUCCESS")
	}
}
