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

type CSharpWriter struct {
	output    string
	namespace string
}

func NewCSharpWriter(dir string) *CSharpWriter {
	var (
		output       = path.Join(dir, "csharp", "Config")
		_, namespace = path.Split(output)
	)

	return &CSharpWriter{
		output:    output,
		namespace: helper.CapitalizeLeading(namespace),
	}
}

func (p *CSharpWriter) mkdir() error {
	return helper.Mkdir(p.OutputDir())
}

func (p *CSharpWriter) OutputDir() string {
	return p.output
}

func (p *CSharpWriter) Write(configMetas []*meta.Config) error {
	if err := p.mkdir(); err != nil {
		return err
	}

	fmt.Println("> write cs configs ...")
	if err := p.writeConfigs(configMetas); err != nil {
		return err
	}
	fmt.Println("< write cs configs SUCCEED !")

	if err := p.writeConfigMgr(configMetas); err != nil {
		return err
	}
	fmt.Println("< write cs ConfigMgr SUCCEED !")

	if err := p.writeEnums(); err != nil {
		return err
	}
	fmt.Println("< write cs enums SUCCEED !")

	return nil
}

func (p *CSharpWriter) writeConfigs(configMetas []*meta.Config) error {
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
			err = p.writeConfig(meta, constsTmpl)
		} else if consts.SideServer(meta.KeyField.Side) {
			err = p.writeConfig(meta, configTmpl)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *CSharpWriter) writeConfig(configMeta *meta.Config, tmpl *gotemplate.Template) error {
	filename := helper.UnderlineToCamelCase(configMeta.Filename, true)
	csFilepath := path.Join(p.OutputDir(), filename+".cs")
	file, err := os.Create(csFilepath)
	if err != nil {
		return err
	}

	conf := p.parseConfig(configMeta)
	return tmpl.Execute(file, &conf)
}

func (p *CSharpWriter) parseConfig(configMeta *meta.Config) template.CSharpConfig {
	conf := template.CSharpConfig{
		Filename:   configMeta.Filename,
		ConfigName: helper.UnderlineToCamelCase(configMeta.Filename, true),
		Namespace:  p.namespace,
	}

	if !configMeta.IsConst {
		conf.KeyType = string(configMeta.KeyField.Type)
		conf.KeyFieldName = p.toFieldName(configMeta.KeyField.Name)
	}

	for _, f := range configMeta.Fields {
		conf.ConfigFields = append(conf.ConfigFields, template.CSharpConfigField{
			Name: p.toFieldName(f.Name),
			Type: p.toTypeName(f.Type, f.TypeParams...),
			Desc: f.Desc,
		})
	}

	return conf
}

func (p *CSharpWriter) toFieldName(fieldName string) string {
	return strings.ToUpper(string(fieldName[0])) + fieldName[1:]
}

func (p *CSharpWriter) toTypeName(dataType consts.DataType, params ...string) string {
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
		return fmt.Sprintf("%s[]", p.toTypeName(elemDataType, elemParams...))
	case consts.Array2:
		elemDataType, elemParams := helper.ParseDataType(params[0])
		return fmt.Sprintf("%s[][]", p.toTypeName(elemDataType, elemParams...))
	case consts.Map:
		keyDataType, keyParams := helper.ParseDataType(params[0])
		keyType := p.toTypeName(keyDataType, keyParams...)
		valDataType, valParams := helper.ParseDataType(params[1])
		valType := p.toTypeName(valDataType, valParams...)
		return fmt.Sprintf("Dictionary<%s, %s>", keyType, valType)
	}

	return string(dataType)
}

func (p *CSharpWriter) writeConfigMgr(metas []*meta.Config) error {
	tmpl, err := gotemplate.New("CSharpConfigMgr").Parse(template.CSharpConfigMgrTemplate)
	if err != nil {
		return err
	}

	conf := template.CSharpConfigMgr{
		Namespace: p.namespace,
		Configs:   []string{},
	}

	for _, meta := range metas {
		configName := helper.UnderlineToCamelCase(meta.Filename, true)
		if !meta.IsConst {
			configName += "Configs"
		}
		conf.Configs = append(conf.Configs, configName)
	}

	csFilepath := path.Join(p.OutputDir(), "ConfigMgr.cs")
	file, err := os.Create(csFilepath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	return tmpl.Execute(file, &conf)
}

func (p *CSharpWriter) writeEnums() error {
	tmpl, err := gotemplate.New("CSharpEnums").Parse(template.CSharpEnumsTemplate)
	if err != nil {
		return err
	}

	enums := types.Enums()
	conf := template.CSharpEnums{
		Namespace: p.namespace,
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

	csFilepath := path.Join(p.OutputDir(), "Enums.cs")
	file, err := os.Create(csFilepath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	return tmpl.Execute(file, &conf)
}
