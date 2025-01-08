package example

import (
	"encoding/json"
	"os"
	"path"
)

var Consts = &constsConfig{}

type constsConfig struct {
	ConstInt    int            `json:"ConstInt"` // 这是一个常量
	ConstInt8   int8           `json:"ConstInt8"`
	ConstInt16  int16          `json:"ConstInt16"`
	ConstInt32  int32          `json:"ConstInt32"`
	ConstInt64  int64          `json:"ConstInt64"`
	ConstUint   uint           `json:"ConstUint"`
	ConstUint8  uint8          `json:"ConstUint8"`
	ConstUint16 uint16         `json:"ConstUint16"`
	ConstUint32 uint32         `json:"ConstUint32"`
	ConstUint64 uint64         `json:"ConstUint64"`
	ConstFloat  float32        `json:"ConstFloat"`
	ConstDouble float64        `json:"ConstDouble"`
	ConstString string         `json:"ConstString"`
	ConstBool   bool           `json:"ConstBool"`
	ConstArray  []string       `json:"ConstArray"`
	ConstArray2 [][]int64      `json:"ConstArray2"`
	ConstMap    map[string]int `json:"ConstMap"`
	ConstEnum   int32          `json:"ConstEnum"` // 周三呀

}

func (c *constsConfig) Load(dir string) {
	data, err := os.ReadFile(path.Join(dir, c.Filename()+".json"))
	if err != nil {
		panic("load config error: " + err.Error())
	}

	if err = json.Unmarshal(data, &c); err != nil {
		panic("parse config error: " + err.Error())
	}
}

func (c *constsConfig) Filename() string {
	return "consts"
}
