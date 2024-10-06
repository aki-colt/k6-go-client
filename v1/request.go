package k6cli

import (
	"encoding/json"
	"fmt"
)

type K6Request interface {
	genRequestScript(imports map[string]any) string
}

type K6HttpRequest struct {
	// name should be distinct
	Name   string
	Method string
	Url    string
	Body   string
	Param  K6HttpRequestParam
	Checks []K6Check
}

type K6HttpRequestParam struct {
	Auth    string
	Cookies map[string]struct {
		Value   string `json:"value"`
		Replace bool   `json:"replace"`
	}
	Headers   map[string]string
	Redirects int64
	Tags      map[string]string
	Timeout   string
}

func (k K6HttpRequest) Type() string {
	return "http"
}

func (k K6HttpRequest) genRequestScript(imports map[string]any) string {
	imports["import http from 'k6/http';"] = nil
	// deal with request
	requestParams := make([]string, 0)
	if k.Param.Auth != "" {
		requestParams = append(requestParams, fmt.Sprintf("\t"+`auth: "%s"`, k.Param.Auth))
	}
	if k.Param.Redirects != 0 {
		requestParams = append(requestParams, fmt.Sprintf("\t"+`redirects: %d`, k.Param.Redirects))
	}
	if k.Param.Timeout != "" {
		requestParams = append(requestParams, fmt.Sprintf("\t"+`timeout: "%s"`, k.Param.Timeout))
	}
	if len(k.Param.Cookies) != 0 {
		// cookieStr := "cookies: {\n"
		// for key, cookie := range k.Param.Cookies {
		// 	cookieStr += fmt.Sprintf("%s: {value: %s, replace: %v}, \n", key, cookie.Value, cookie.Replace)
		// }
		// cookieStr += "}"
		cookie, _ := json.Marshal(k.Param.Cookies)
		cookieStr := fmt.Sprintf("\t"+"cookies: %s", cookie)
		requestParams = append(requestParams, cookieStr)
	}
	if len(k.Param.Headers) != 0 {
		// headerStr := "headers: {\n"
		// for key, v := range k.Param.Headers {
		// 	headerStr += fmt.Sprintf("%s: %s,\n", key, v)
		// }
		// headerStr += "}"
		headers, _ := json.Marshal(k.Param.Headers)
		headerStr := fmt.Sprintf("\t"+"headers: %s", headers)
		requestParams = append(requestParams, headerStr)
	}
	if len(k.Param.Tags) != 0 {
		// tagStr := "tags: {\n"
		// for key, v := range k.Param.Tags {
		// 	tagStr += fmt.Sprintf("%s: %s,\n", key, v)
		// }
		// tagStr += "}"
		tags, _ := json.Marshal(k.Param.Tags)
		tagStr := fmt.Sprintf("\t"+"tags: %s", tags)
		requestParams = append(requestParams, tagStr)
	}
	param := "null"
	if len(requestParams) != 0 {
		param = "{\n"
		for _, v := range requestParams {
			param += v + ",\n"
		}
		param += "}"
	}
	body := "null"
	if k.Body != "" {
		body = `"` + k.Body + `"`
	}
	// deal with check
	checks := ""
	if len(k.Checks) != 0 {
		imports["import { check } from 'k6';"] = nil
		checks = fmt.Sprintf("check(%s, {\n", k.Name)
		for _, check := range k.Checks {
			checks += check.genCheckScript()
		}
		checks += "})"
	}
	// generate res
	res := fmt.Sprintf("let %s = http.request(\"%s\", \"%s\", %s, %s);\n%s;",
		k.Name, k.Method, k.Url, body, param, checks)
	return res
}
