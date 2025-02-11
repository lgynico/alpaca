package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lgynico/alpaca/meta"
)

type (
	Tempaltes struct {
		Config    string
		ConfigMgr string
		Consts    string
		Enums     string
	}

	FileWriter interface {
		OutputDir() string
		Write(configMetas []*meta.Config) error
	}

	NoneWriter struct {
	}
)

func (n *NoneWriter) OutputDir() string {
	return ""
}

func (n *NoneWriter) Write(configMetas []*meta.Config) error {
	fmt.Println("Do nothing")
	return nil
}

func readTemplates(lang, tmplPath string) (Tempaltes, error) {
	config, err := os.ReadFile(filepath.Join(tmplPath, lang, "config.tmpl"))
	if err != nil {
		return Tempaltes{}, err
	}

	configMgr, err := os.ReadFile(filepath.Join(tmplPath, lang, "config_mgr.tmpl"))
	if err != nil {
		return Tempaltes{}, err
	}

	consts, err := os.ReadFile(filepath.Join(tmplPath, lang, "consts.tmpl"))
	if err != nil {
		return Tempaltes{}, err
	}

	enums, err := os.ReadFile(filepath.Join(tmplPath, lang, "enums.tmpl"))
	if err != nil {
		return Tempaltes{}, err
	}

	return Tempaltes{
		Config:    string(config),
		ConfigMgr: string(configMgr),
		Consts:    string(consts),
		Enums:     string(enums),
	}, nil
}
