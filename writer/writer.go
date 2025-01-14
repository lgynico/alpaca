package writer

import (
	"fmt"

	"github.com/lgynico/alpaca/meta"
)

type FileWriter interface {
	OutputDir() string
	Write(configMetas []*meta.Config) error
}

type NoneWriter struct {
}

func (n *NoneWriter) OutputDir() string {
	return ""
}

func (n *NoneWriter) Write(configMetas []*meta.Config) error {
	fmt.Println("Do nothing")
	return nil
}
