package meta

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/types"
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

		if strings.HasPrefix(entry.Name(), types.EnumFilename) ||
			strings.HasPrefix(entry.Name(), "__const__") {
			continue
		}

		meta, err := parse(path.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}

		fmt.Printf("parse [%s] SUCCEED !\r\n", entry.Name())
		metaList = append(metaList, meta...)
	}

	fmt.Printf("parse excel files SUCCEED (total: %d)\r\n", len(metaList))
	return metaList, nil
}

func parse(filepath string) ([]*Config, error) {
	file, err := excelize.OpenFile(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	prefix := prefix(filepath)
	metas := []*Config{}

	for _, sheet := range file.WorkBook.Sheets.Sheet {
		rows, err := file.GetRows(sheet.Name)
		if err != nil {
			return nil, err
		}

		if len(rows) < 2 {
			return nil, errors.New("row count lack")
		}

		rows = fixedRows(rows)

		var (
			nameRow []string
			typeRow []string
			sideRow []string
			descRow []string
			ruleRow []string
		)

		i := 0
		for ; i < len(rows); i++ {
			row := rows[i]
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
			}
		}

		if nameRow == nil || typeRow == nil {
			return nil, errors.New("name/type row lack")
		}

		if sideRow == nil {
			sideRow = make([]string, len(rows[0]))
		}
		if descRow == nil {
			descRow = make([]string, len(rows[0]))
		}
		if ruleRow == nil {
			ruleRow = make([]string, len(rows[0]))
		}

		meta := &Config{
			Filename: prefix + "_" + strings.ToLower(sheet.Name),
			Fields:   []*Field{},
		}

		for j := 1; j < len(nameRow); j++ {
			dataType, typeParams := consts.ParseDataType(typeRow[j])
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

		for ; i < len(rows); i++ {
			for j := 1; j < len(rows[i]); j++ {
				meta.Fields[j-1].RawValues = append(meta.Fields[j-1].RawValues, rows[i][j])
			}
		}

		metas = append(metas, meta)
	}

	return metas, nil
}

func prefix(filepath string) string {
	filename := path.Base(filepath)
	return strings.TrimSuffix(filename, path.Ext(filename))
}

func fixedRows(rows [][]string) [][]string {
	rowLen := len(rows[0])
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) == rowLen {
			continue
		}

		rows[i] = append(rows[i], make([]string, rowLen-len(rows[i]))...)
	}

	return rows
}
