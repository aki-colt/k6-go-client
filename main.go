package main

import (
	"log"

	k6cli "github.com/aki-colt/k6-go-client/v1"
)

func main() {
	bInner := k6cli.K6NormalBlock{
		Constants: map[string]any{
			"con1": 1,
			"con2": "aaa",
		},
		Variables: map[string]any{
			"v1": []string{"1", "2"},
			"v2": map[string]any{
				"v2a": "a",
				"v2b": 13,
			},
		},
	}
	client := k6cli.K6Client{
		Options: k6cli.K6Options{
			Vus:        10,
			Iterations: 20,
		},
		Imports: []string{},
		Block: k6cli.K6NormalBlock{
			Constants: map[string]any{
				"con1": 1,
				"con2": "aaa",
			},
			Variables: map[string]any{
				"v1": []string{"1", "2"},
				"v2": map[string]any{
					"v2a": "a",
					"v2b": 13,
				},
			},
			Requests: []k6cli.K6Request{
				k6cli.K6HttpRequest{
					Name:   "test1",
					Method: "GET",
					Url:    "https://test.k6.io",
					Param: k6cli.K6HttpRequestParam{
						Timeout:   "10s",
						Redirects: 10,
						Headers: map[string]string{
							"X-MyHeader": "k6test",
						},
						Cookies: map[string]struct {
							Value   string `json:"value"`
							Replace bool   `json:"replace"`
						}{
							"my-cookie": {"my-value", true},
						},
						Tags: map[string]string{
							"my-tag": "tag1",
						},
					},
					Checks: []k6cli.K6Check{
						{Name: "my-check", Check: "(res) => res.status === 200"},
					},
				},
			},
			Blocks: []k6cli.K6Block{
				k6cli.K6IfBlock{
					Condition: "con1 === 1",
					IfBlock:   bInner,
				},
				k6cli.K6ForBlock{
					Condition: "let i = 0; i < 2; i++",
					Block:     bInner,
				},
			},
		},
	}
	res, err := client.GenScript()
	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Print(res)
	}
}
