package writer

import (
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/meta"
)

func WriteJSON(filepath string, configMeta *meta.Config) error {
	jsonStr := toJSON(configMeta)
	jsonFilepath := path.Join(filepath, configMeta.Filename+".json")
	return os.WriteFile(jsonFilepath, []byte(jsonStr), os.ModePerm)
}

func toJSON(configMeta *meta.Config) string {
	jsonArr := make([]string, len(configMeta.Fields[0].Values))

	for i, field := range configMeta.Fields {
		for j := 0; j < len(field.Values); j++ {
			if i > 0 {
				jsonArr[j] += ","
			}

			value := consts.ToString(reflect.ValueOf(field.Values[j]))
			jsonArr[j] += fmt.Sprintf("%q:%s", field.Name, value)
		}
	}

	jsonStr := "[\r\n"
	for i, v := range jsonArr {
		jsonStr += "\t{" + v + "}"
		if i < len(jsonArr)-1 {
			jsonStr += ","
		}
		jsonStr += "\r\n"
	}
	jsonStr += "]"

	return jsonStr
}
