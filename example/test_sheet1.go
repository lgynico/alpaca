package example

import (
	"encoding/json"
	"os"
	"path"
)

var TestSheet1 = &testSheet1Config{}

type (
	testSheet1Row struct {
		I       int            `json:"i"` // 一个int类型值
		I32     int32          `json:"i32"`
		I64     int64          `json:"i64"`
		U       uint           `json:"u"`
		U8      uint8          `json:"u8"`
		U32     uint32         `json:"u32"`
		U64     uint64         `json:"u64"`
		F       float32        `json:"f"`
		D       float64        `json:"d"`
		S       string         `json:"s"`
		A       []int          `json:"a"`
		A2      [][]string     `json:"a2"`
		M       map[int]string `json:"m"`
		B       bool           `json:"b"`
		Weekday int32          `json:"weekday"`
	}

	testSheet1Config struct {
		rows map[int]*testSheet1Row
	}
)

func (c *testSheet1Config) Load(dir string) {
	data, err := os.ReadFile(path.Join(dir, c.Filename()+".json"))
	if err != nil {
		panic("load config error: " + err.Error())
	}

	rows := []*testSheet1Row{}
	if err = json.Unmarshal(data, &rows); err != nil {
		panic("parse config error: " + err.Error())
	}

	c.rows = map[int]*testSheet1Row{}
	for _, row := range rows {
		c.rows[row.I] = row
	}

}

func (c *testSheet1Config) Filename() string {
	return "test_sheet1"
}

func (c *testSheet1Config) Get(key int) (*testSheet1Row, bool) {
	row, ok := c.rows[key]
	return row, ok
}

func (c *testSheet1Config) List() []*testSheet1Row {
	list := []*testSheet1Row{}
	for _, row := range c.rows {
		list = append(list, row)
	}
	return list
}
