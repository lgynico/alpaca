package template

type Placeholder string

const (
	PackageName      Placeholder = "${packageName}"
	StructName       Placeholder = "${structName}"
	ConfigName       Placeholder = "${configName}"
	ExportConfigName Placeholder = "${exportConfigName}"
	KeyType          Placeholder = "${keyType}"
	KeyFieldName     Placeholder = "${keyFieldName}"
	Filename         Placeholder = "${filename}"
	Fields           Placeholder = "${fields}"
	RegisterConfigs  Placeholder = "${registerConfigs}"
)
