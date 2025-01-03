package writer

import "github.com/lgynico/alpaca/meta"

type WriteFunc func(filepath string, configMate *meta.Config) error
