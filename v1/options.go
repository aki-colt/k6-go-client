package k6cli

import "fmt"

type K6Options struct {
	Vus        int64
	Duration   string
	Iterations int64
}

func (o K6Options) genOptionScript() string {
	res := "export const options = {\n"
	if o.Vus != 0 {
		res += fmt.Sprintf("vus: %d,\n", o.Vus)
	}
	if o.Duration != "" {
		res += fmt.Sprintf("duration: %s,\n", o.Duration)
	}
	if o.Iterations != 0 {
		res += fmt.Sprintf("iterations: %d,\n", o.Iterations)
	}
	res += "}"
	return res
}
