package template

import _ "embed"

var (
	//go:embed templates/golang/config.tmpl
	GoConfigTemplate string

	//go:embed templates/golang/config_mgr.tmpl
	GoConfigMgrTemplate string

	//go:embed templates/golang/enums.tmpl
	GoEnumsTemplate string

	//go:embed templates/golang/consts.tmpl
	GoConstsTemplate string
)

type (
	GoConfig struct {
		Package      string
		Filename     string
		ConfigName   string
		ExportName   string
		RowName      string
		RowFields    []string
		KeyType      string
		KeyFieldName string
	}

	GoConfigMgr struct {
		Package string
		Configs []string
	}

	GoEnum struct {
		Key   string
		Value int32
	}

	GoEnums struct {
		Package string
		Enums   [][]GoEnum
	}
)
