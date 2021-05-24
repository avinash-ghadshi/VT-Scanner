package config

import (
	"os"

	"github.com/seebs/gogetopt"
)

var API_KEY = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
var DEBUG = false

var VTAPIS = map[string]string{
	"FILE_SUBMITION_URL": "https://www.virustotal.com/api/v3/files",
	"FILE_ANALYSIS_URL":  "https://www.virustotal.com/api/v3/analyses/%s",
}

type UserOptions struct {
	File   string
	Url    string
	Apikey string
}

func (uo *UserOptions) GetInputs() (true bool) {
	opts, _, _ := gogetopt.GetOpt(os.Args[1:], "f:u:k:")

	for k, v := range opts {
		if k == "f" {
			uo.File = v.Value
		} else if k == "u" {
			uo.Url = v.Value
		} else if k == "k" {
			uo.Apikey = v.Value
		}
	}
	if uo.Apikey == "" {
		uo.Apikey = API_KEY
	}
	return
}
