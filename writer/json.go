package writer

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/helper"
	"github.com/lgynico/alpaca/meta"
)

type JsonWriter struct {
	output string
	c, s   bool
}

func NewJsonWriter(dir string, c, s bool) *JsonWriter {
	return &JsonWriter{output: path.Join(dir, "json"), c: c, s: s}
}

func (p *JsonWriter) mkdir() error {
	dir := []string{}
	if p.c {
		dir = append(dir, path.Join(p.OutputDir(), consts.OutputClient))
	}

	if p.s {
		dir = append(dir, path.Join(p.OutputDir(), consts.OutputServer))
	}

	if len(dir) == 0 {
		return nil
	}

	return helper.Mkdir(dir...)
}

func (p *JsonWriter) OutputDir() string {
	return p.output
}

func (p *JsonWriter) Write(configMetas []*meta.Config) error {
	if err := p.mkdir(); err != nil {
		return err
	}

	fmt.Println("> write json ...")
	for _, m := range configMetas {
		if err := p.write(m); err != nil {
			return err
		}
		fmt.Printf("write [%s.json] SUCCEED !\r\n", m.Filename)
	}
	fmt.Println("< write json SUCCEED !")
	return nil
}

func (p *JsonWriter) write(configMeta *meta.Config) error {
	var (
		jsonStr      string
		jsonFilepath string
	)

	if p.s {
		if configMeta.IsConst || consts.SideServer(configMeta.KeyField.Side) {
			jsonStr = p.stringify(configMeta, consts.SideServer)
			jsonFilepath = path.Join(p.OutputDir(), consts.OutputServer, configMeta.Filename+".json")
			if err := os.WriteFile(jsonFilepath, []byte(jsonStr), os.ModePerm); err != nil {
				return err
			}
		}
	}

	if p.c {
		if configMeta.IsConst || consts.SideClient(configMeta.KeyField.Side) {
			jsonStr = p.stringify(configMeta, consts.SideClient)
			jsonFilepath = path.Join(p.OutputDir(), consts.OutputClient, configMeta.Filename+".json")
			if err := os.WriteFile(jsonFilepath, []byte(jsonStr), os.ModePerm); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *JsonWriter) stringify(configMeta *meta.Config, side consts.Side) string {
	jsonArr := make([]string, len(configMeta.Fields[0].Values))

	for i, field := range configMeta.Fields {

		if !side(field.Side) {
			continue
		}

		for j := 0; j < len(field.Values); j++ {
			if configMeta.IsConst {
				jsonArr[j] += "\t"
			}

			var value string
			if field.Values[j] != nil {
				value = helper.ToString(reflect.ValueOf(field.Values[j]))
			} else {
				value = helper.DefaultStringValue(field.Type)
			}

			jsonArr[j] += fmt.Sprintf("%q:%s", field.Name, value)

			if i < len(configMeta.Fields)-1 {
				jsonArr[j] += ","
			}

			if configMeta.IsConst {
				jsonArr[j] += "\r\n"
			}
		}
	}

	if configMeta.IsConst {
		return fmt.Sprintf("{\r\n%s}", jsonArr[0])
	}

	jsonStr := "[\r\n"
	for i, v := range jsonArr {
		jsonStr += "\t{" + strings.TrimRight(v, ",") + "}"
		if i < len(jsonArr)-1 {
			jsonStr += ","
		}
		jsonStr += "\r\n"
	}
	jsonStr += "]"

	return jsonStr
}
