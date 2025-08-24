package ali

import (
	config "cloud_firewall/config"
	"encoding/json"
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	swasopen "github.com/alibabacloud-go/swas-open-20200601/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
	"strings"
)

// CreateClient Description:
//
// 使用凭据初始化账号Client
//
// @return Client
//
// @throws Exception
func CreateClient() (_result *swasopen.Client, _err error) {
	// 工程代码建议使用更安全的无AK方式，凭据配置方式请参见：https://help.aliyun.com/document_detail/378661.html。
	c, _err := credential.NewCredential(&credential.Config{
		Type:            tea.String("access_key"),
		AccessKeyId:     tea.String(config.Cfg.Ali.Key),
		AccessKeySecret: tea.String(config.Cfg.Ali.Secret),
	})
	if _err != nil {
		return _result, _err
	}

	cfg := &openapi.Config{
		Credential: c,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/SWAS-OPEN
	cfg.Endpoint = tea.String("swas.cn-shanghai.aliyuncs.com")
	_result = &swasopen.Client{}
	_result, _err = swasopen.NewClient(cfg)
	return _result, _err
}

func CreateFirewallRules(request *swasopen.CreateFirewallRulesRequest) error {
	client, err := CreateClient()
	if err != nil {
		return err
	}
	return catchAliError(func() error {
		_, _err := client.CreateFirewallRules(request)
		return _err
	})
}

func ModifyFirewallRules(request *swasopen.ModifyFirewallRuleRequest) error {
	client, err := CreateClient()
	if err != nil {
		return err
	}
	return catchAliError(func() error {
		_, _err := client.ModifyFirewallRule(request)
		return _err
	})
}

func GetFirewallRules(request *swasopen.ListFirewallRulesRequest) (rules []*swasopen.ListFirewallRulesResponseBodyFirewallRules, _err error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}
	_err = catchAliError(func() error {
		result, err := client.ListFirewallRules(request)
		if err != nil {
			return err
		}
		if result.Body != nil {
			rules = result.Body.GetFirewallRules()
		}
		return nil
	})
	return
}

func checkErr(tryErr error) error {
	if tryErr == nil {
		return nil
	}
	var err = &tea.SDKError{}
	if !errors.As(tryErr, &err) {
		return tryErr
	}
	var data interface{}
	d := json.NewDecoder(strings.NewReader(tea.StringValue(err.Data)))
	_ = d.Decode(&data)
	if m, ok := data.(map[string]interface{}); ok {
		recommend, _ := m["Recommend"]
		recommendText, _ := util.AssertAsString(recommend)
		if recommendText != nil {
			return errors.New(*recommendText)
		}
	}
	msg, _ := util.AssertAsString(err.Message)
	if msg != nil {
		return errors.New(*msg)
	}
	return tryErr
}

func catchAliError(f func() error) (err error) {
	defer func() {
		if r := tea.Recover(recover()); r != nil {
			err = r
		}
	}()
	err = checkErr(f())
	return
}
