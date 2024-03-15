package wecom

import (
	"encoding/json"
	"fmt"

	"github.com/hkloudou/nw"
	"github.com/tidwall/gjson"
)

type third_Licence_order struct {
	OrderID   string `json:"order_id"`
	OrderType int    `json:"order_type"`
}

type third_Licence_order_accout struct {
	ActiveCode string `json:"active_code"`
	UserID     string `json:"userid"`
	Type       int    `json:"type"`
}

type third_licence_actived_account struct {
	UserID     string `json:"userid"`
	Type       int    `json:"type"`
	ExpireTime int    `json:"expire_time"`
	ActiveTime int    `json:"active_time"`
}

type third_licence_acticecode_info struct {
	UserID     string `json:"userid"`
	Type       int    `json:"type"`
	Status     int    `json:"status"`
	ExpireTime int    `json:"expire_time"`
	ActiveTime int    `json:"active_time"`
}

// "active_code": "code1",
// 		"type": 1,
// 		"status": 1,
// 		"userid": "USERID",
// 		"create_time":1640966400,
// 		"active_time": 1640966400,
// 		"expire_time":1640966400,
// 		"merge_info":
// 		{
// 			  "to_active_code":"code_new",
// 			  "from_active_code":"code_old"
// 		},
// 		"share_info":
// 		{
// 			"to_corpid":"CORPID",
// 			"from_corpid":"CORPID"
// 		}

func Third_Licence_list_order(provider_access_token string, corpid string) ([]third_Licence_order, error) {
	return postPaged[third_Licence_order](
		"order_list",
		map[string]any{
			"corpid": corpid,
			"limit":  1000,
		}, nw.WithSite(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/license/list_order?provider_access_token=%s", provider_access_token)),
	)
}

func Third_Licence_list_order_account(provider_access_token string, orderid string) ([]third_Licence_order_accout, error) {
	return postPaged[third_Licence_order_accout](
		"account_list",
		map[string]any{
			"order_id": orderid,
			"limit":    1000,
		}, nw.WithSite(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/license/list_order_account?provider_access_token=%s", provider_access_token)),
	)
}

// 账号管理 - 激活账号
func Third_License_list_actived_account(provider_access_token string, corpid string) ([]third_licence_actived_account, error) {
	return postPaged[third_licence_actived_account](
		"account_list",
		map[string]any{
			"corpid": corpid,
			"limit":  1000,
		},
		nw.WithSite(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/license/list_actived_account?provider_access_token=%s", provider_access_token)),
		// nw.WithLog(true),
	)
}

// 账号管理 - 获得激活码详情
func Third_License_get_active_info_by_code(provider_access_token string, corpid string, activecode string) (*third_licence_acticecode_info, error) {
	return nw.PostJsonData[third_licence_acticecode_info](map[string]any{
		"active_code": activecode,
		"corpid":      corpid,
	},
		nw.WithSite(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/license/get_active_info_by_code?provider_access_token=%s", provider_access_token)),
		nw.WithDataKeys("active_info"),
	)
	// return postPaged[third_licence_actived_account](
	// 	"account_list",
	// 	map[string]any{
	// 		"corpid": corpid,
	// 		"limit":  1000,
	// 	},
	// 	nw.WithSite(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/license/list_actived_account?provider_access_token=%s", provider_access_token)),
	// 	nw.WithLog(true),
	// )
}

func Third_Licence_active_account(provider_access_token string, corpid string, userid string, active_code string) error {
	_, err := nw.PostJsonData[gjson.Result](map[string]any{
		"active_code": active_code,
		"corpid":      corpid,
		"userid":      userid,
	},
		nw.WithSite(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/license/active_account?provider_access_token=%s", provider_access_token)),
	)
	return err
}

//https://qyapi.weixin.qq.com/cgi-bin/license/active_account?provider_access_token=ACCESS_TOKEN

func postPaged[T any](dataKey string, requestData map[string]any, opts ...nw.NwOption) ([]T, error) {
	var next = ""
	var results = make([]T, 0)
	i := 0
	for {
		i++
		if i > 4 {
			break
		}
		if next != "" {
			requestData["cursor"] = next
		}
		opts = append(opts, nw.WithDataKeys(""))
		g, err := nw.PostJsonData[gjson.Result](requestData,
			opts...,
		)
		if err != nil {
			return nil, err
		}
		if !g.Get(dataKey).Exists() {
			return nil, fmt.Errorf("not exist")
		}
		var response []T
		err = json.Unmarshal([]byte(g.Get(dataKey).Raw), &response)
		if err != nil {
			return nil, err
		}
		results = append(results, response...)
		if g.Get("has_more").Int() == 1 {
			next = g.Get("next_cursor").String()
		}
		if next == "" {
			break
		}
	}
	return results, nil
}
