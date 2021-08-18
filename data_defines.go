package main

import "encoding/json"

type Result struct {
	status int
	body   interface{}
}

func (p *Result) Json() string {
	data, err := json.Marshal(p)
	if err != nil {
		data, _ = json.Marshal(Result{status: -1})
		return string(data)
	}
	return string(data)
}
