package template

import _ "embed"

var (
	//go:embed templates/csharp/config.tmpl
	CSharpConfigTemplate string

	//go:embed templates/csharp/config_mgr.tmpl
	CSharpConfigMgrTemplate string

	//go:embed templates/csharp/enums.tmpl
	CSharpEnumsTemplate string

	//go:embed templates/csharp/consts.tmpl
	CSharpConstsTemplate string
)

type (
	CSharpConfigField struct {
		Type string
		Name string
		Desc string
	}

	CSharpConfig struct {
		Namespace    string
		Filename     string
		ConfigName   string
		ConfigFields []CSharpConfigField
		KeyType      string
		KeyFieldName string
	}

	CSharpConfigMgr struct {
		Namespace string
		Configs   []string
	}

	CSharpEnums struct {
		Namespace string
		Enums     []CSharpEnum
	}

	CSharpEnum struct {
		Name   string
		Fields [][]string // [[key, value]]
	}
)
