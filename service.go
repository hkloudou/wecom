package wecom

import (
	"fmt"

	"github.com/hkloudou/nw"
)

func Third_Service_get_provider_token(providerCorpID string, providerSecret string) (*string, error) {
	return nw.PostJsonData[string](map[string]any{
		"corpid":          providerCorpID,
		"provider_secret": providerSecret,
	},
		nw.WithSite("https://qyapi.weixin.qq.com/cgi-bin/service/get_provider_token"),
		nw.WithDataKeys("provider_access_token"),
	)
}

func Third_Service_corpid_to_opencorpid(provider_access_token string, corpid string) (*string, error) {
	return nw.PostJsonData[string](map[string]any{
		"corpid": corpid,
	},
		nw.WithSite(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/service/corpid_to_opencorpid?provider_access_token=%s", provider_access_token)),
		nw.WithDataKeys("open_corpid"),
	)
}

type SuiteToken struct {
	SuiteAccessToken string `json:"suite_access_token"`
	ExpiresIn        int    `json:"expires_in"`
}

func Third_Service_get_suite_token(suite_id, suite_secret, suite_ticket string) (*SuiteToken, error) {
	return nw.PostJsonData[SuiteToken](map[string]any{
		"suite_id":     suite_id,
		"suite_secret": suite_secret,
		"suite_ticket": suite_ticket,
	},
		nw.WithSite("https://qyapi.weixin.qq.com/cgi-bin/service/get_suite_token"),
		nw.WithDataKeys(""),
	)
}

// func Third_Service_get_admin_list(suite_access_token string, auth_corpid string, agentid int) (*string, error) {
// 	return nw.PostJsonData[string](map[string]any{
// 		"corpid": corpid,
// 	},
// 		nw.WithSite(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/service/corpid_to_opencorpid?provider_access_token=%s", provider_access_token)),
// 		nw.WithDataKeys("open_corpid"),
// 	)
// }

// https://qyapi.weixin.qq.com/cgi-bin/service/get_admin_list?suite_access_token=SUITE_ACCESS_TOKEN
