package main

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/hkloudou/wecom"
	"github.com/hkloudou/xlib/xcolor"
)

//go:embed corps.txt
var str string

func main() {

	corpids := strings.Split(str, "\n")
	for i := 0; i < len(corpids); i++ {
		corpid := corpids[i]
		if len(corpid) < 10 {
			continue
		}
		list, err := wecom.Third_License_list_actived_account(provider_access_token, corpid)
		if err != nil {
			fmt.Println(err)
			continue
		}
		actived := []string{}
		for i := 0; i < len(list); i++ {
			//
			if list[i].Type == 2 && list[i].ExpireTime > int(time.Now().Add(time.Hour*24*19).Unix()) {
				// fmt.Println(list[i].UserID, xcolor.Green("已激活"))
				actived = append(actived, list[i].UserID)
				continue
			}
			// fmt.Println(list[i].ExpireTime)
			// panic("x")
			if strings.Contains(strings.Join(actived, ","), list[i].UserID) {
				// fmt.Println(list[i].UserID, xcolor.Green("已激活2"))
				continue
			}
			activecode, err := getLastAviabeCode(provider_access_token, corpid)

			if err != nil || len(activecode) < 5 {
				fmt.Println(xcolor.Red(corpid), "没有足够的Token")
				continue
			}
			err = wecom.Third_Licence_active_account(provider_access_token, corpid, list[i].UserID, activecode)
			if err != nil {
				continue
			}
			fmt.Println(xcolor.Green("激活成功"))
		}
	}
}

func getLastAviabeCode(provider_access_token string, corpid string) (string, error) {
	orders, err := wecom.Third_Licence_list_order(provider_access_token, corpid)
	if err != nil {
		return "", err
	}
	for i := 0; i < len(orders); i++ {
		if orders[i].OrderType != 1 {
			continue
		}
		codes,
			err := wecom.Third_Licence_list_order_account(provider_access_token, orders[i].OrderID)
		if err != nil {
			return "", err
		}
		for j := 0; j < len(codes); j++ {
			if len(codes[j].UserID) > 0 {
				continue
			}
			if codes[j].Type != 2 {
				continue
			}
			info, err := wecom.Third_License_get_active_info_by_code(provider_access_token, corpid, codes[j].ActiveCode)
			if err != nil {
				return "", err
			}
			if info.Status != 1 {
				// fmt.Println(xcolor.Yellow(codes[j].ActiveCode), info.Status)
				continue
			}
			// fmt.Println("codes[i]", codes[i].Type, codes[i].UserID, codes[i].ActiveCode, codes[i])
			return codes[j].ActiveCode, nil
		}
	}
	return "", fmt.Errorf("找不到")
}
