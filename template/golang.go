package template

const GoConfig = `
package ${packageName}

import (
	"os"
	"encoding/json"
	"path"
)

var ${exportConfigName} = &${configName}Config{}

type (
	${structName}Row struct {
		${fields}
	}

	${configName}Config struct {
		rows map[${keyType}]*${structName}Row
	}
)

func (c *${configName}Config) Load(dir string) {
	data, err := os.ReadFile(path.Join(dir, c.Filename() + ".json"))
	if err != nil {
		panic("load config error: " + err.Error())
	}

	rows := []*${structName}Row{}
	if err = json.Unmarshal(data, &rows); err != nil {
		panic("parse config error: " + err.Error())
	}

	c.rows = map[${keyType}]*${structName}Row{}
	for _, row := range rows {
		c.rows[row.${keyFieldName}] = row
	}

}

func (c *${configName}Config) Filename() string {
	return "${filename}"
}

func (c *${configName}Config) Get(key ${keyType}) (*${structName}Row, bool) {
	row, ok := c.rows[key]
	return row, ok
}

func (c *${configName}Config) List() []*${structName}Row {
	list := []*${structName}Row{}
	for _, row := range c.rows {
		list = append(list, row)
	}
	return list
}
`

const GoField = "${exportFieldName} ${filedType} `json:\"${fieldName}\"`"

const GoConfigMgr = `
package ${packageName}

import (
	"fmt"
)

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
	${registerConfigs}
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
`

const GoRegister = `c.Register(${exportConfigName})`
