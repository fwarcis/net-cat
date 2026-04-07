package flags

import (
	"encoding"
	"flag"
)

type VarType interface {
	encoding.TextUnmarshaler
	encoding.TextMarshaler
}

func TextVar[T VarType](name string, value T, usage string) (p T) {
	p = func(v T) T {
		return v
	}(value)
	flag.TextVar(p, name, value, usage)
	return p
}

func IsSet(name string) bool {
	res := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			res = true
			return
		}
	})
	return res
}
