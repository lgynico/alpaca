package writer

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/lgynico/alpaca/helper"
	"github.com/lgynico/alpaca/types"

	gotemplate "text/template"

	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/meta"
	"github.com/lgynico/alpaca/template"
)

var (
	//go:embed template/golang/config.tmpl
	GoConfigTemplate string

	//go:embed template/golang/config_mgr.tmpl
	GoConfigMgrTemplate string

	//go:embed template/golang/enums.tmpl
	GoEnumsTemplate string

	//go:embed template/golang/consts.tmpl
	GoConstsTemplate string
)

func WriteGoConfigs(filepath string, configMetas []*meta.Config) error {
	configTmpl, err := gotemplate.New("GoConfig").Parse(GoConfigTemplate)
	if err != nil {
		return err
	}
	constsTmpl, err := gotemplate.New("GoConsts").Parse(GoConstsTemplate)
	if err != nil {
		return err
	}

	for _, meta := range configMetas {
		if meta.IsConst {
			err = writeGoConfig(meta, constsTmpl, filepath)
		} else {
			err = writeGoConfig(meta, configTmpl, filepath)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func writeGoConfig(configMeta *meta.Config, tmpl *gotemplate.Template, filepath string) error {
	goFilepath := path.Join(filepath, configMeta.Filename+".go")
	file, err := os.Create(goFilepath)
	if err != nil {
		return err
	}

	conf := parseGoConfig(configMeta)
	_, pkgName := path.Split(filepath)

	conf.Package = pkgName

	return tmpl.Execute(file, &conf)
}

func parseGoConfig(configMeta *meta.Config) template.GoConfig {
	var (
		filename     = configMeta.Filename
		configName   = helper.UnderlineToCamelCase(configMeta.Filename, false)
		exportName   = helper.UnderlineToCamelCase(configMeta.Filename, true)
		rowName      = configName
		keyType      = consts.Unknown
		keyFieldName = ""
		fields       []string
	)

	if !configMeta.IsConst {
		keyType = configMeta.KeyField.Type
		keyFieldName = toGoFieldName(configMeta.KeyField.Name, true)
	}

	for _, f := range configMeta.Fields {
		var (
			fieldName = toGoFieldName(f.Name, true)
			goType    = toGoType(f.Type, f.TypeParams...)
			field     string
		)

		if f.Desc == "" {
			field = fmt.Sprintf("%s %s `json:\"%s\"`\r\n", fieldName, goType, f.Name)
		} else {
			field = fmt.Sprintf("%s %s `json:\"%s\"` // %s\r\n", fieldName, goType, f.Name, f.Desc)
		}

		fields = append(fields, field)
	}

	return template.GoConfig{
		Filename:     filename,
		ConfigName:   configName,
		ExportName:   exportName,
		RowName:      rowName,
		RowFields:    fields,
		KeyType:      string(keyType),
		KeyFieldName: keyFieldName,
	}
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
		elemDataType, elemParams := helper.ParseDataType(params[0])
		return fmt.Sprintf("[]%s", toGoType(elemDataType, elemParams...))
	case consts.Array2:
		elemDataType, elemParams := helper.ParseDataType(params[0])
		return fmt.Sprintf("[][]%s", toGoType(elemDataType, elemParams...))
	case consts.Map:
		keyDataType, keyParams := helper.ParseDataType(params[0])
		keyType := toGoType(keyDataType, keyParams...)
		valDataType, valParams := helper.ParseDataType(params[1])
		valType := toGoType(valDataType, valParams...)
		return fmt.Sprintf("map[%s]%s", keyType, valType)
	case consts.Enum:
		return "int32"
	}

	return string(dataType)
}

func WriteGoConfigMgr(filepath string, metas []*meta.Config) error {
	tmpl, err := gotemplate.New("GoConfigMgr").Parse(GoConfigMgrTemplate)
	if err != nil {
		return err
	}

	_, pkgName := path.Split(filepath)
	conf := template.GoConfigMgr{
		Package: pkgName,
		Configs: []string{},
	}

	for _, meta := range metas {
		exportName := helper.UnderlineToCamelCase(meta.Filename, true)
		conf.Configs = append(conf.Configs, exportName)
	}

	goFilepath := path.Join(filepath, "config_mgr.go")
	file, err := os.Create(goFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, &conf)
}

func WriteGoEnums(filepath string, enums []*types.EnumType) error {
	tmpl, err := gotemplate.New("GoEnums").Parse(GoEnumsTemplate)
	if err != nil {
		return err
	}

	_, pkgName := path.Split(filepath)
	conf := template.GoEnums{
		Package: pkgName,
		Enums:   make([][]template.GoEnum, 0, len(enums)),
	}

	for _, enumType := range enums {
		var goEnum []template.GoEnum
		for _, node := range enumType.Nodes {
			name := fmt.Sprintf("%s_%s", enumType.Name, node.Key)
			goEnum = append(goEnum, template.GoEnum{
				Key:   name,
				Value: node.Value,
			})
		}
		conf.Enums = append(conf.Enums, goEnum)
	}

	goFilepath := path.Join(filepath, "enums.go")
	file, err := os.Create(goFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, &conf)
}
