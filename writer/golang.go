package writer

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/mate"
	"github.com/lgynico/alpaca/template"
)

func WriteGoConfig(filepath string, configMeta *mate.Config) error {
	_, pkgName := path.Split(filepath)
	goStr := toGoConfig(configMeta, pkgName)
	goFilepath := path.Join(filepath, configMeta.Filename+".go")

	return os.WriteFile(goFilepath, []byte(goStr), os.ModePerm)
}

func toGoConfig(configMeta *mate.Config, pkgName string) string {
	var (
		structName       = consts.UnderlineToCamelCase(configMeta.Filename, true)
		configName       = consts.UnderlineToCamelCase(configMeta.Filename, false)
		exportConfigName = structName
		keyType          = configMeta.KeyField.Type
		keyFieldName     = toGoFieldName(configMeta.KeyField.Name, true)
		filename         = configMeta.Filename
		fields           = ""
	)

	for _, f := range configMeta.Fields {
		var (
			fieldName = toGoFieldName(f.Name, true)
			goType    = toGoType(f.Type, f.TypeParams...)
		)

		if f.Desc == "" {
			fields += fmt.Sprintf("%s %s `json:\"%s\"`\r\n", fieldName, goType, f.Name)
		} else {
			fields += fmt.Sprintf("%s %s `json:\"%s\"` // %s\r\n", fieldName, goType, f.Name, f.Desc)
		}
	}

	goStr := strings.ReplaceAll(template.GoConfig, string(template.StructName), structName)
	goStr = strings.ReplaceAll(goStr, string(template.ConfigName), configName)
	goStr = strings.ReplaceAll(goStr, string(template.ExportConfigName), exportConfigName)
	goStr = strings.ReplaceAll(goStr, string(template.KeyType), string(keyType))
	goStr = strings.ReplaceAll(goStr, string(template.KeyFieldName), keyFieldName)
	goStr = strings.ReplaceAll(goStr, string(template.Fields), fields)
	goStr = strings.ReplaceAll(goStr, string(template.Filename), filename)
	goStr = strings.ReplaceAll(goStr, string(template.PackageName), pkgName)
	return goStr
}

func toGoFieldName(fieldName string, export bool) string {
	if !export {
		return fieldName
	}

	return strings.ToUpper(string(fieldName[0])) + fieldName[1:]
}

func toGoType(dataType consts.DataType, params ...string) string {
	switch dataType {
	case consts.Float:
		return "float32"
	case consts.Double:
		return "float64"
	case consts.Array:
		elemDataType, elemParams := consts.ParseDataType(params[0])
		return fmt.Sprintf("[]%s", toGoType(elemDataType, elemParams...))
	case consts.Array2:
		elemDataType, elemParams := consts.ParseDataType(params[0])
		return fmt.Sprintf("[][]%s", toGoType(elemDataType, elemParams...))
	case consts.Map:
		keyDataType, keyParams := consts.ParseDataType(params[0])
		keyType := toGoType(keyDataType, keyParams...)
		valDataType, valParams := consts.ParseDataType(params[1])
		valType := toGoType(valDataType, valParams...)
		return fmt.Sprintf("map[%s]%s", keyType, valType)
	case consts.Enum:
		return "int32"
	}

	return string(dataType)
}

func WriteGoConfigMgr(filepath string, mates []*mate.Config) error {
	_, pkgName := path.Split(filepath)

	registerConfigs := ""
	for _, mate := range mates {
		exportConfigName := consts.UnderlineToCamelCase(mate.Filename, true)
		registerConfigs += strings.ReplaceAll(template.GoRegister, string(template.ExportConfigName), exportConfigName)
		registerConfigs += "\n"
	}

	goStr := strings.ReplaceAll(template.GoConfigMgr, string(template.PackageName), pkgName)
	goStr = strings.ReplaceAll(goStr, string(template.RegisterConfigs), registerConfigs)

	goFilepath := path.Join(filepath, "configmgr.go")

	return os.WriteFile(goFilepath, []byte(goStr), os.ModePerm)
}

