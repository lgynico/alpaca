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

func WriteJSON(filepath string, configMeta *meta.Config) error {
	var (
		jsonStr      string
		jsonFilepath string
	)

	jsonStr = toJSON(configMeta, consts.SideServer)
	jsonFilepath = path.Join(filepath, consts.OutputServer, configMeta.Filename+".json")
	if err := os.WriteFile(jsonFilepath, []byte(jsonStr), os.ModePerm); err != nil {
		return err
	}

	jsonStr = toJSON(configMeta, consts.SideClient)
	jsonFilepath = path.Join(filepath, consts.OutputClient, configMeta.Filename+".json")
	return os.WriteFile(jsonFilepath, []byte(jsonStr), os.ModePerm)
}

func toJSON(configMeta *meta.Config, side consts.Side) string {
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
