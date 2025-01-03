package meta

import "github.com/lgynico/alpaca/consts"

type Field struct {
	Name       string          // 字段名
	Type       consts.DataType // 字段类型
	TypeParams []string        // 类型的参数
	Desc       string          // 字段描述
	Side       string          // 生成端
	Rule       string          // 验证规则
	RuleMeta   Rule            // 验证规则
	RawValues  []string        // 原始值列表
	Values     []any           // 解析值列表
}
