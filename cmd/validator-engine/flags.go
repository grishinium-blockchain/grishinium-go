package main

import (
	"strings"
)

// multiFlag collects repeated string flags into a slice.
type multiFlag []string

func (m *multiFlag) String() string {
	if m == nil {
		return ""
	}
	return strings.Join(*m, ",")
}

func (m *multiFlag) Set(s string) error {
	*m = append(*m, s)
	return nil
}
