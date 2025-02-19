package meta

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/helper"
	"github.com/xuri/excelize/v2"
)

func Parse(dir string) ([]*Config, error) {
	fmt.Println("parse excel files ... ")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	metaList := []*Config{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if !(ext == ".xlsx" || ext == ".xls") {
			continue
		}

		if strings.HasPrefix(entry.Name(), consts.FilenameEnum) {
			continue
		}

		isConst := strings.HasPrefix(entry.Name(), consts.FilenameConst)
		meta, err := parse(path.Join(dir, entry.Name()), isConst)
		if err != nil {
			return nil, err
		}

		fmt.Printf("parse [%s] SUCCEED !\r\n", entry.Name())
		metaList = append(metaList, meta...)
	}

	fmt.Printf("parse excel files SUCCEED (total: %d)\r\n", len(metaList))
	return metaList, nil
}

func parse(filepath string, isConst bool) ([]*Config, error) {
	file, err := excelize.OpenFile(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	prefix := prefix(filepath)
	metas := []*Config{}

	for _, sheet := range file.WorkBook.Sheets.Sheet {
		var (
			datas [][]string
			err   error
		)

		if isConst {
			datas, err = file.GetCols(sheet.Name)
		} else {
			datas, err = file.GetRows(sheet.Name)
		}

		if err != nil {
			return nil, err
		}

		if len(datas) < 2 {
			return nil, errors.New("row count lack")
		}

		datas = fixedRows(datas, filepath, sheet.Name)

		var (
			nameRow  []string
			typeRow  []string
			sideRow  []string
			descRow  []string
			ruleRow  []string
			valueRow []string // only for const
		)

		i := 0
		for ; i < len(datas); i++ {
			row := datas[i]
			if len(row[0]) == 0 {
				break
			}

			switch consts.RowType(row[0]) {
			case consts.NameRow:
				nameRow = row
			case consts.TypeRow:
				typeRow = row
			case consts.SideRow:
				sideRow = row
			case consts.DescRow:
				descRow = row
			case consts.RuleRow:
				ruleRow = row
			case consts.ValueRow:
				valueRow = row
			}
		}

		if nameRow == nil || typeRow == nil {
			return nil, errors.New("name/type row lack")
		}
		if isConst && valueRow == nil {
			return nil, errors.New("lack value row for const")
		}

		if sideRow == nil {
			sideRow = make([]string, len(datas[0]))
		}
		if descRow == nil {
			descRow = make([]string, len(datas[0]))
		}
		if ruleRow == nil {
			ruleRow = make([]string, len(datas[0]))
		}

		var filename string
		if isConst {
			filename = "consts"
		} else {
			filename = prefix + "_" + strings.ToLower(sheet.Name)
		}

		meta := &Config{
			IsConst:  isConst,
			Filename: filename,
			Fields:   []*Field{},
		}

		for j := 1; j < len(nameRow); j++ {
			// skip empty row
			if len(nameRow[j]) == 0 {
				continue
			}

			dataType, typeParams := helper.ParseDataType(typeRow[j])
			field := &Field{
				Name:       nameRow[j],
				Type:       dataType,
				TypeParams: typeParams,
				Desc:       descRow[j],
				Side:       sideRow[j],
				Rule:       ruleRow[j],
				RawValues:  []string{},
			}
			meta.Fields = append(meta.Fields, field)
		}

		if isConst {
			for i := 0; i < len(meta.Fields); i++ {
				meta.Fields[i].RawValues = append(meta.Fields[i].RawValues, valueRow[i+1])
			}
		} else {
			for ; i < len(datas); i++ {
				k := 0
				for j := 1; j < len(nameRow); j++ {
					if len(nameRow[j]) == 0 {
						continue
					}
					meta.Fields[k].RawValues = append(meta.Fields[k].RawValues, datas[i][j])
					k++
				}
			}
		}

		metas = append(metas, meta)
	}

	if isConst {
		metaMerge := metas[0]
		for i := 1; i < len(metas); i++ {
			metaMerge.Fields = append(metaMerge.Fields, metas[i].Fields...)
		}
		metas = []*Config{metaMerge}
	}

	return metas, nil
}

func prefix(filepath string) string {
	filename := path.Base(filepath)
	return strings.TrimSuffix(filename, path.Ext(filename))
}

func fixedRows(rows [][]string, filename, sheetName string) [][]string {
	rowLen := len(rows[0])
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) < rowLen {
			rows[i] = append(rows[i], make([]string, rowLen-len(rows[i]))...)
		} else if len(rows[i]) > rowLen {
			fmt.Printf("\033[33m[WARN] %s:%s row %d contains extra space\033[0m\r\n", filename, sheetName, i+1)
		}
	}

	return rows
}

// func removeRowValue(rows []string, idx int) []string {
// 	if len(rows) <= idx {
// 		return rows
// 	}

// 	return append(rows[:idx], rows[idx+1:]...)
// }
