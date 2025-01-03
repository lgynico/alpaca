package types

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type (
	EnumNode struct {
		Key   string
		Value int32
	}

	EnumType struct {
		Name     string
		HasValue bool
		Nodes    []*EnumNode
	}

	enumRows struct {
		keyRow   []string
		valueRow []string
	}
)

var enums = map[string]*EnumType{}

func GetEnum(name string) (*EnumType, bool) {
	enum, ok := enums[name]
	return enum, ok
}

func Enums() map[string]*EnumType {
	return enums
}

func ParseEnum(dir string) error {
	var (
		filename = fmt.Sprintf("%s/%s.xlsx", dir, EnumFilename)
		file     *excelize.File
		err      error
	)

	if file, err = excelize.OpenFile(filename); err != nil {
		filename = fmt.Sprintf("%s/%s.xls", dir, EnumFilename)
		if file, err = excelize.OpenFile(filename); err != nil {
			return nil
		}
	}

	for _, sheet := range file.WorkBook.Sheets.Sheet {
		rows, err := file.GetRows(sheet.Name)
		if err != nil {
			return err
		}

		enumRows := groupEnumRows(rows)
		for name, enumRow := range enumRows {
			enum := &EnumType{
				Name:     name,
				HasValue: enumRow.valueRow != nil,
			}

			for i, key := range enumRow.keyRow {
				value := int32(i)
				if enum.HasValue {
					v, err := strconv.ParseInt(enumRow.valueRow[i], 10, 32)
					if err != nil {
						return err
					}
					value = int32(v)
				}

				enum.Nodes = append(enum.Nodes, &EnumNode{
					Key:   key,
					Value: value,
				})
			}

			enums[name] = enum
		}
	}

	return nil
}

func groupEnumRows(rows [][]string) map[string]*enumRows {
	m := make(map[string]*enumRows)
	var currEnum *enumRows
	for _, row := range rows {
		if len(row) < 2 {
			continue
		}

		if len(row[0]) > 1 {
			currEnum = &enumRows{}
			m[row[0]] = currEnum
		}

		if row[1] == EnumFieldKey {
			currEnum.keyRow = row[2:]
		} else if row[1] == EnumFieldValue {
			currEnum.valueRow = row[2:]
		}
	}
	return m
}
