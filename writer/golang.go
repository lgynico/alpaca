package writer

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/lgynico/alpaca/helper"
	"github.com/lgynico/alpaca/types"

	gotemplate "text/template"

	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/meta"
	"github.com/lgynico/alpaca/template"
)

type GoWriter struct {
	output  string
	pkgName string
}

func NewGoWriter(dir string) *GoWriter {
	var (
		output     = path.Join(dir, "go", "conf")
		_, pkgName = path.Split(output)
	)

	return &GoWriter{
		output:  path.Join(dir, "go", "conf"),
		pkgName: pkgName,
	}
}

func (p *GoWriter) mkdir() error {
	return helper.Mkdir(p.OutputDir())
}

func (p *GoWriter) OutputDir() string {
	return p.output
}

func (p *GoWriter) Write(configMetas []*meta.Config) error {
	if err := p.mkdir(); err != nil {
		return err
	}

	fmt.Println("> write go configs ...")
	if err := p.writeConfigs(configMetas); err != nil {
		return err
	}
	fmt.Println("< write go configs SUCCEED !")

	if err := p.writeConfigMgr(configMetas); err != nil {
		return err
	}
	fmt.Println("< write go ConfigMgr SUCCEED !")

	if err := p.writeEnums(); err != nil {
		return err
	}
	fmt.Println("< write go enums SUCCEED !")

	p.formatCodes()

	return nil
}

func (p *GoWriter) writeConfigs(configMetas []*meta.Config) error {
	configTmpl, err := gotemplate.New("GoConfig").Parse(template.GoConfigTemplate)
	if err != nil {
		return err
	}
	constsTmpl, err := gotemplate.New("GoConsts").Parse(template.GoConstsTemplate)
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

func (p *GoWriter) writeConfig(configMeta *meta.Config, tmpl *gotemplate.Template) error {
	goFilepath := path.Join(p.OutputDir(), configMeta.Filename+".go")
	file, err := os.Create(goFilepath)
	if err != nil {
		return err
	}

	conf := p.parseConfig(configMeta)
	return tmpl.Execute(file, &conf)
}

func (p *GoWriter) parseConfig(configMeta *meta.Config) template.GoConfig {
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
		keyFieldName = p.toFieldName(configMeta.KeyField.Name, true)
	}

	for _, f := range configMeta.Fields {
		if !consts.SideServer(f.Side) {
			continue
		}

		var (
			fieldName = p.toFieldName(f.Name, true)
			goType    = p.toTypeName(f.Type, f.TypeParams...)
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
		Package:      p.pkgName,
		Filename:     filename,
		ConfigName:   configName,
		ExportName:   exportName,
		RowName:      rowName,
		RowFields:    fields,
		KeyType:      string(keyType),
		KeyFieldName: keyFieldName,
	}
}

func (p *GoWriter) toFieldName(fieldName string, export bool) string {
	if !export {
		return fieldName
	}

	return strings.ToUpper(string(fieldName[0])) + fieldName[1:]
}

func (p *GoWriter) toTypeName(dataType consts.DataType, params ...string) string {
	switch dataType {
	case consts.Float:
		return "float32"
	case consts.Double:
		return "float64"
	case consts.Array:
		elemDataType, elemParams := helper.ParseDataType(params[0])
		return fmt.Sprintf("[]%s", p.toTypeName(elemDataType, elemParams...))
	case consts.Array2:
		elemDataType, elemParams := helper.ParseDataType(params[0])
		return fmt.Sprintf("[][]%s", p.toTypeName(elemDataType, elemParams...))
	case consts.Map:
		keyDataType, keyParams := helper.ParseDataType(params[0])
		keyType := p.toTypeName(keyDataType, keyParams...)
		valDataType, valParams := helper.ParseDataType(params[1])
		valType := p.toTypeName(valDataType, valParams...)
		return fmt.Sprintf("map[%s]%s", keyType, valType)
	case consts.Enum:
		return "int32"
	}

	return string(dataType)
}

func (p *GoWriter) writeConfigMgr(metas []*meta.Config) error {
	tmpl, err := gotemplate.New("GoConfigMgr").Parse(template.GoConfigMgrTemplate)
	if err != nil {
		return err
	}

	conf := template.GoConfigMgr{
		Package: p.pkgName,
		Configs: []string{},
	}

	for _, m := range metas {
		if m.IsConst || consts.SideServer(m.KeyField.Side) {
			exportName := helper.UnderlineToCamelCase(m.Filename, true)
			conf.Configs = append(conf.Configs, exportName)
		}
	}

	goFilepath := path.Join(p.OutputDir(), "config_mgr.go")
	file, err := os.Create(goFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, &conf)
}

func (p *GoWriter) writeEnums() error {
	tmpl, err := gotemplate.New("GoEnums").Parse(template.GoEnumsTemplate)
	if err != nil {
		return err
	}

	enums := types.Enums()
	conf := template.GoEnums{
		Package: p.pkgName,
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

	goFilepath := path.Join(p.OutputDir(), "enums.go")
	file, err := os.Create(goFilepath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	return tmpl.Execute(file, &conf)
}

func (p *GoWriter) formatCodes() {
	if err := os.Chdir(p.OutputDir()); err == nil {
		if err := exec.Command("gofmt", "-w", ".").Run(); err != nil {
			fmt.Printf("format go codes FAILED: %v\r\n", err)
		} else {
			fmt.Println("< format go codes SUCCEED !")
		}
	} else {
		fmt.Printf("format go codes FAILED: %v\r\n", err)
	}
}
