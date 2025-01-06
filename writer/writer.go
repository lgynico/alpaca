package writer

import "github.com/lgynico/alpaca/meta"

type WriteFunc func(filepath string, configMeta *meta.Config) error
