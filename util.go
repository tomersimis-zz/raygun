package raygun

import (
	"strings"

	"github.com/kaeuferportal/stack2struct"
)

func GetCurrentStack() StackTrace {
	s := StackTrace{}
	stack2struct.Current(&s)

	offset := 1 // package github.com/kaeuferportal/stack2struct
	for strings.HasPrefix(s[offset].PackageName, PackageName) {
		offset++
	}

	return s[offset:]
}
