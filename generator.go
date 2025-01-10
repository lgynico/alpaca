package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/helper"
	"github.com/lgynico/alpaca/meta"
	"github.com/lgynico/alpaca/types"
	"github.com/lgynico/alpaca/writer"
)

var generators = []func(metas []*meta.Config) error{
	genJSON,
	genGO,
	genCSharp,
}

func genJSON(metas []*meta.Config) error {
	if len(json_out) == 0 {
		return nil
	}

	if err := helper.Mkdir(
		path.Join(json_out, consts.OutputClient),
		path.Join(json_out, consts.OutputServer),
	); err != nil {
		return err
	}

	fmt.Println("> write json ...")
	for _, meta := range metas {
		if err := writer.WriteJSON(json_out, meta); err != nil {
			return err
		}
		fmt.Printf("write [%s.json] SUCCEED !\r\n", meta.Filename)
	}
	fmt.Println("< write json SUCCEED !")

	return nil
}

func genGO(metas []*meta.Config) error {
	if len(go_out) == 0 {
		return nil
	}

	if err := helper.Mkdir(go_out); err != nil {
		return err
	}

	fmt.Println("> write go configs ...")
	if err := writer.WriteGoConfigs(go_out, metas); err != nil {
		return err
	}
	fmt.Println("< write go configs SUCCEED !")

	if err := writer.WriteGoConfigMgr(go_out, metas); err != nil {
		return err
	}
	fmt.Println("< write go ConfigMgr SUCCEED !")

	if err := writer.WriteGoEnums(go_out, types.Enums()); err != nil {
		return err
	}
	fmt.Println("< write go enums SUCCEED !")

	if err := os.Chdir(go_out); err == nil {
		if err := exec.Command("gofmt", "-w", ".").Run(); err != nil {
			fmt.Printf("format go codes FAILED: %v\r\n", err)
		} else {
			fmt.Println("< format go codes SUCCEED !")
		}
	} else {
		fmt.Printf("format go codes FAILED: %v\r\n", err)
	}

	return nil
}

func genCSharp(metas []*meta.Config) error {
	if len(cs_out) == 0 {
		return nil
	}

	if err := helper.Mkdir(cs_out); err != nil {
		return err
	}

	fmt.Println("> write cs configs ...")
	if err := writer.WriteCSharpConfigs(cs_out, metas); err != nil {
		return err
	}
	fmt.Println("< write cs configs SUCCEED !")

	if err := writer.WriteCSharpConfigMgr(cs_out, metas); err != nil {
		return err
	}
	fmt.Println("< write cs ConfigMgr SUCCEED !")

	if err := writer.WriteCSharpEnums(cs_out, types.Enums()); err != nil {
		return err
	}
	fmt.Println("< write cs enums SUCCEED !")

	return nil
}
