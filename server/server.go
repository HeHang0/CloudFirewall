package server

import (
	"cloud_firewall/ali"
	"cloud_firewall/config"
	"encoding/json"
	swasopen "github.com/alibabacloud-go/swas-open-20200601/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type RequestAddData struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Type     string `json:"type"`
	Token    string `json:"token"`
	Region   string `json:"region"`
	Remark   string `json:"remark"`
	Message  string `json:"message"`
	Protocol string `json:"protocol"`
	Instance string `json:"instance"`
}

const TypeUpdate = "update"

// 判断是否为公网 IP
func isPublicIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return false
	}

	// IPv4 私有地址范围
	privateIPBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"100.64.0.0/10", // CGNAT
	}

	for _, cidr := range privateIPBlocks {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return false
		}
	}

	// IPv6 本地地址范围
	privateIPv6Blocks := []string{
		"fc00::/7",  // Unique local address
		"fe80::/10", // Link-local
	}

	for _, cidr := range privateIPv6Blocks {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return false
		}
	}

	return true
}

func getClientIP(r *http.Request) string {
	// 优先 X-Forwarded-For
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// 可能有多个 IP，用逗号分隔，第一个才是客户端
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}

	// 其次 X-Real-IP
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 最后 RemoteAddr
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func getAddBody(r *http.Request) (*RequestAddData, int) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed
	}

	body, err := io.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		log.Printf("read body err: %v", err)
		return nil, http.StatusBadRequest
	}

	if len(body) == 0 {
		return nil, http.StatusLengthRequired
	}
	var data RequestAddData
	if err = json.Unmarshal(body, &data); err != nil {
		data.Message = err.Error()
		return &data, http.StatusBadRequest
	}

	if data.Token != config.Cfg.Token {
		return nil, http.StatusUnauthorized
	}

	if data.IP == "" {
		data.IP = getClientIP(r)
	}
	data.IP = strings.TrimSpace(data.IP)
	ip := net.ParseIP(data.IP)
	if ip == nil {
		data.Message = "不是合法IP！"
		return &data, http.StatusBadRequest
	}
	if !isPublicIP(ip) {
		data.Message = "请填写IP参数！"
		return &data, http.StatusBadRequest
	}
	return &data, http.StatusOK
}

func findExistsRule(rules []*swasopen.ListFirewallRulesResponseBodyFirewallRules, data *RequestAddData) *swasopen.ListFirewallRulesResponseBodyFirewallRules {
	update := data.Type == TypeUpdate
	for i := 0; i < len(rules); i++ {
		rule := rules[i]
		condition := false
		if update {
			condition = rule.Remark != nil && *rule.Remark == data.Remark
		} else {
			condition = rule.SourceCidrIp != nil && *rule.SourceCidrIp == data.IP
		}
		if condition && rule.RuleId != nil && rule.RuleProtocol != nil &&
			strings.ToLower(*rule.RuleProtocol) == strings.ToLower(data.Protocol) &&
			rule.Port != nil && *rule.Port == strconv.Itoa(data.Port) {
			return rule
		}
	}
	return nil
}

func handlerAdd(w http.ResponseWriter, r *http.Request) {
	data, code := getAddBody(r)
	if code != http.StatusOK {
		msg := http.StatusText(code)
		if data != nil {
			msg = data.Message
		}
		http.Error(w, msg, code)
		return
	}
	rules, err := ali.GetFirewallRules(&swasopen.ListFirewallRulesRequest{
		RegionId:   tea.String(data.Region),
		InstanceId: tea.String(data.Instance),
		PageSize:   tea.Int32(100),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rule := findExistsRule(rules, data)

	if rule != nil && (data.Type != TypeUpdate || *rule.SourceCidrIp == data.IP) {
		http.Error(w, "已存在！", http.StatusOK)
		return
	}
	var message string
	if rule != nil {
		err = ali.ModifyFirewallRules(&swasopen.ModifyFirewallRuleRequest{
			RuleId:       rule.RuleId,
			Port:         tea.String(strconv.Itoa(data.Port)),
			Remark:       tea.String(data.Remark),
			RegionId:     tea.String(data.Region),
			InstanceId:   tea.String(data.Instance),
			RuleProtocol: tea.String(data.Protocol),
			SourceCidrIp: tea.String(data.IP),
		})
		message = "更新成功！"
	} else {
		err = ali.CreateFirewallRules(&swasopen.CreateFirewallRulesRequest{
			RegionId:   tea.String(data.Region),
			InstanceId: tea.String(data.Instance),
			FirewallRules: []*swasopen.CreateFirewallRulesRequestFirewallRules{
				{
					Port:         tea.String(strconv.Itoa(data.Port)),
					Remark:       tea.String(data.Remark),
					RuleProtocol: tea.String(data.Protocol),
					SourceCidrIp: tea.String(data.IP),
				},
			},
		})
		message = "添加成功!"
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.Error(w, message, http.StatusOK)
	}
}

func init() {
	http.HandleFunc("/ali/add", handlerAdd)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Hello World!", http.StatusOK)
	})
}
