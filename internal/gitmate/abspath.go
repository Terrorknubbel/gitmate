package gitmate

type AbsPath string

func (a AbsPath) Empty() bool {
	return a == ""
}
