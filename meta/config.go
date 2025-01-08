package meta

type Config struct {
	Filename string   // 文件名
	Fields   []*Field // 字段
	KeyField *Field   // 主键字段
	IsConst  bool
}
