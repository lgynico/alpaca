package consts

type DataType string

const (
	Unknown  DataType = ""
	Bool     DataType = "bool"
	Int      DataType = "int"
	Int8     DataType = "int8"
	Int16    DataType = "int16"
	Int32    DataType = "int32"
	Int64    DataType = "int64"
	Uint     DataType = "uint"
	Uint8    DataType = "uint8"
	Uint16   DataType = "uint16"
	Uint32   DataType = "uint32"
	Uint64   DataType = "uint64"
	Float    DataType = "float"
	Double   DataType = "double"
	String   DataType = "string"
	Array    DataType = "array"
	Array2   DataType = "array2"
	Map      DataType = "map"
	Enum     DataType = "enum"
	Datetime DataType = "datetime"
)
