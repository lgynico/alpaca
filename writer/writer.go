package writer

import "github.com/lgynico/alpaca/mate"

type WriteFunc func(filepath string, configMate *mate.Config) error
