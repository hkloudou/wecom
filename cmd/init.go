package main

import (
	_ "embed"

	"github.com/hkloudou/wecom"
)

//go:embed pcorpid.txt
var provider_corpid string

//go:embed psecret.txt
var provider_secret string

var provider_access_token string

func init() {
	tmp, err := wecom.Third_Service_get_provider_token(provider_corpid, provider_secret)
	if err != nil {
		panic(err)
	}
	provider_access_token = *tmp
}
