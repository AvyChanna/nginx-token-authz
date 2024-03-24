package tabletest

import (
	"testing"
)

type T[Fields any, Arg any] []struct {
	Name   string
	Fields Fields
	Args   []Arg
}

func (table T[Fields, Arg]) Run(t2 *testing.T, f func(*testing.T, Fields, []Arg)) {
	for _, tCase := range table {
		t2.Run(tCase.Name, func(t *testing.T) {
			f(t, tCase.Fields, tCase.Args)
		})
	}
}
