package mate

type Config struct {
	Filename string   // 文件名
	Fields   []*Field // 字段
	KeyField *Field   // 主键字段
}
