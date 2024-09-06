package gitmate

import (
	"path/filepath"
)

type AbsPath string

func (a AbsPath) Empty() bool {
	return a == ""
}

func (a AbsPath) IsAbs() bool {
	return filepath.IsAbs(a.String())
}

func (a AbsPath) String() string {
	return string(a)
}
