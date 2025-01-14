package main

import (
	"github.com/lgynico/alpaca/meta"
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

	w := writer.NewJsonWriter(json_out)
	return w.Write(metas)
}

func genGO(metas []*meta.Config) error {
	if len(go_out) == 0 {
		return nil
	}

	w := writer.NewGoWriter(go_out)
	return w.Write(metas)
}

func genCSharp(metas []*meta.Config) error {
	if len(cs_out) == 0 {
		return nil
	}

	w := writer.NewCSharpWriter(cs_out)
	return w.Write(metas)
}
