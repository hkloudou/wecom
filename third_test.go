package wecom_test

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"

	"github.com/hkloudou/wecom"
	"github.com/hkloudou/xlib/xcolor"
)

//go:embed cmd/pcorpid.txt
var provider_corpid string

//go:embed cmd/psecret.txt
var provider_secret string

var provider_access_token string

var corpid = "wpWCvICQAAR0XB7rpTvvefcQAuX6vGig"

func init() {
	tmp, err := wecom.Third_Service_get_provider_token(provider_corpid, provider_secret)
	if err != nil {
		panic(err)
	}
	provider_access_token = *tmp
}

func Test_ThirdLicenceOrderList(t *testing.T) {
	list, err := wecom.Third_Licence_list_order(provider_access_token, "wpWCvICQAApZUhkWKRlzo7txkj2871og")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}

func Test_ThirdLicenceListOrderAccount(t *testing.T) {
	list, err := wecom.Third_Licence_list_order_account(provider_access_token, "OI000009C82B5863124DFB7DF7CABT")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)

}

// func Test

func Test_Third_License_list_actived_account(t *testing.T) {
	list, err := wecom.Third_License_list_actived_account(provider_access_token, corpid)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}

func Test_autoActive(t *testing.T) {
	// t.Log(getLastAviabeCode(provider_access_token, "wpWCvICQAApZUhkWKRlzo7txkj2871og"))
	list, err := wecom.Third_License_list_actived_account(provider_access_token, corpid)
	if err != nil {
		t.Fatal(err)
	}
	actived := []string{}
	for i := 0; i < len(list); i++ {
		if list[i].Type == 2 {
			fmt.Println(list[i].UserID, xcolor.Green("已激活"))
			actived = append(actived, list[i].UserID)
			continue
		}
		if strings.Contains(strings.Join(actived, ","), list[i].UserID) {
			fmt.Println(list[i].UserID, xcolor.Green("已激活2"))
			continue
		}
		activecode, err := getLastAviabeCode(provider_access_token, corpid)

		if err != nil {
			t.Fatal(err, provider_access_token, corpid)
		}
		err = wecom.Third_Licence_active_account(provider_access_token, corpid, list[i].UserID, activecode)
		if err != nil {
			continue
			t.Fatal(err, provider_access_token, corpid, list[i].UserID, activecode)
		}
		fmt.Println(xcolor.Green("激活成功"))
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
