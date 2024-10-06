package k6cli

import "fmt"

type K6Check struct {
	// name should be distinct
	Name string
	// check condition in js grammar
	// Example:
	// (res) => res.status === 200
	Check string
}

func (k K6Check) genCheckScript() string {
	return fmt.Sprintf("\t'%s': %s,\n", k.Name, k.Check)
}
