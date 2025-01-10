package writer

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/lgynico/alpaca/helper"
	"github.com/lgynico/alpaca/types"

	gotemplate "text/template"

	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/meta"
	"github.com/lgynico/alpaca/template"
)

func WriteCSharpConfigs(filepath string, configMetas []*meta.Config) error {
	configTmpl, err := gotemplate.New("CSharpConfig").Parse(template.CSharpConfigTemplate)
	if err != nil {
		return err
	}
	constsTmpl, err := gotemplate.New("CSharpConsts").Parse(template.CSharpConstsTemplate)
	if err != nil {
		return err
	}

	for _, meta := range configMetas {
		if meta.IsConst {
			err = writeCSharpConfig(meta, constsTmpl, filepath)
		} else if consts.SideServer(meta.KeyField.Side) {
			err = writeCSharpConfig(meta, configTmpl, filepath)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func writeCSharpConfig(configMeta *meta.Config, tmpl *gotemplate.Template, filepath string) error {
	filename := helper.UnderlineToCamelCase(configMeta.Filename, true)
	csFilepath := path.Join(filepath, filename+".cs")
	file, err := os.Create(csFilepath)
	if err != nil {
		return err
	}

	conf := parseCSharpConfig(configMeta)
	_, namespace := path.Split(filepath)

	conf.Namespace = helper.CapitalizeLeading(namespace)

	return tmpl.Execute(file, &conf)
}

func parseCSharpConfig(configMeta *meta.Config) template.CSharpConfig {
	conf := template.CSharpConfig{
		Filename:   configMeta.Filename,
		ConfigName: helper.UnderlineToCamelCase(configMeta.Filename, true),
	}

	if !configMeta.IsConst {
		conf.KeyType = string(configMeta.KeyField.Type)
		conf.KeyFieldName = toCSharpFieldName(configMeta.KeyField.Name)
	}

	for _, f := range configMeta.Fields {
		conf.ConfigFields = append(conf.ConfigFields, template.CSharpConfigField{
			Name: toCSharpFieldName(f.Name),
			Type: toCSharpType(f.Type, f.TypeParams...),
			Desc: f.Desc,
		})
	}

	return conf
}

func toCSharpFieldName(fieldName string) string {
	return strings.ToUpper(string(fieldName[0])) + fieldName[1:]
}

func toCSharpType(dataType consts.DataType, params ...string) string {
	switch dataType {
	case consts.Int, consts.Int32, consts.Enum:
		return "int"
	case consts.Int8:
		return "sbyte"
	case consts.Int16:
		return "short"
	case consts.Int64:
		return "long"
	case consts.Uint, consts.Uint32:
		return "uint"
	case consts.Uint8:
		return "byte"
	case consts.Uint16:
		return "ushort"
	case consts.Uint64:
		return "ulong"
	case consts.Float:
		return "float"
	case consts.Double:
		return "double"
	case consts.Array:
		elemDataType, elemParams := helper.ParseDataType(params[0])
		return fmt.Sprintf("%s[]", toCSharpType(elemDataType, elemParams...))
	case consts.Array2:
		elemDataType, elemParams := helper.ParseDataType(params[0])
		return fmt.Sprintf("%s[][]", toCSharpType(elemDataType, elemParams...))
	case consts.Map:
		keyDataType, keyParams := helper.ParseDataType(params[0])
		keyType := toCSharpType(keyDataType, keyParams...)
		valDataType, valParams := helper.ParseDataType(params[1])
		valType := toCSharpType(valDataType, valParams...)
		return fmt.Sprintf("Dictionary<%s, %s>", keyType, valType)
	}

	return string(dataType)
}

func WriteCSharpConfigMgr(filepath string, metas []*meta.Config) error {
	tmpl, err := gotemplate.New("CSharpConfigMgr").Parse(template.CSharpConfigMgrTemplate)
	if err != nil {
		return err
	}

	_, namespace := path.Split(filepath)
	conf := template.CSharpConfigMgr{
		Namespace: helper.CapitalizeLeading(namespace),
		Configs:   []string{},
	}

	for _, meta := range metas {
		configName := helper.UnderlineToCamelCase(meta.Filename, true)
		if !meta.IsConst {
			configName += "Configs"
		}
		conf.Configs = append(conf.Configs, configName)
	}

	csFilepath := path.Join(filepath, "ConfigMgr.cs")
	file, err := os.Create(csFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, &conf)
}

func WriteCSharpEnums(filepath string, enums []*types.EnumType) error {
	tmpl, err := gotemplate.New("CSharpEnums").Parse(template.CSharpEnumsTemplate)
	if err != nil {
		return err
	}

	_, namespace := path.Split(filepath)
	conf := template.CSharpEnums{
		Namespace: helper.CapitalizeLeading(namespace),
		Enums:     make([]template.CSharpEnum, 0, len(enums)),
	}

	for _, enumType := range enums {
		csharpEnum := template.CSharpEnum{
			Name:   enumType.Name,
			Fields: [][]string{},
		}

		for _, node := range enumType.Nodes {
			csharpEnum.Fields = append(csharpEnum.Fields,
				[]string{node.Key, strconv.Itoa(int(node.Value))})
		}
		conf.Enums = append(conf.Enums, csharpEnum)
	}

	csFilepath := path.Join(filepath, "Enums.cs")
	file, err := os.Create(csFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, &conf)
}
