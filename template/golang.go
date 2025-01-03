package template

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
