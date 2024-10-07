# K6-go-client

**K6-go-client** is a light go library to generate k6 script in js by go code with nearly no dependencies. It accepts go struct and output string of k6 js script.

## getting started
It works like bellow in your project
```go
import (
	"fmt"
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
	log.Print(client.GenScript())
}
```
and will get a result string as follow
```js
import http from 'k6/http';
import { check } from 'k6';

export const options = {
vus: 10,
iterations: 20,
}

export default function() {
const con1 = 1
const con2 = "aaa"

let v2 = {"v2a":"a","v2b":13}
let v1 = ["1","2"]

let test1 = http.request("GET", "https://test.k6.io", null, {
	redirects: 10,
	timeout: "10s",
	cookies: {"my-cookie":{"value":"my-value","replace":true}},
	headers: {"X-MyHeader":"k6test"},
	tags: {"my-tag":"tag1"},
});
check(test1, {
	'my-check': (res) => res.status === 200,
});

if (con1 === 1) {
const con1 = 1
const con2 = "aaa"

let v1 = ["1","2"]
let v2 = {"v2a":"a","v2b":13}


} else {

}

for (let i = 0; i < 2; i++) {
const con1 = 1
const con2 = "aaa"

let v1 = ["1","2"]
let v2 = {"v2a":"a","v2b":13}


}

}
```

As you can see, the code indent is not very well, but it works in k6. 

Attention that **the package will only output string**, if you need to feed the script to a k6 programme, you should write the string into a js file and feed it to k6. This should be easy with go.

The package now only accepts http request and very simple options. It will be a feature to join more options and requests like ws and so on. **But it should be only used when you need to generate simple script** because generated script will have only one function and is not friendly for human to read.