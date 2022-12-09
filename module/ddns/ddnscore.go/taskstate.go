package ddnscore

import (
	"time"
)

// 固定的主域名
var staticMainDomains = []string{"com.cn", "org.cn", "net.cn", "ac.cn", "eu.org"}

// 获取ip失败的次数

// Domains Ipv4/Ipv6 DDNSTaskState
type DDNSTaskState struct {
	IpAddr              string
	Domains             []Domain
	domainsRawStrList   []string
	WebhookCallTime     string    `json:"WebhookCallTime"`     //最后触发时间
	WebhookCallResult   bool      `json:"WebhookCallResult"`   //触发结果
	WebhookCallErrorMsg string    `json:"WebhookCallErrorMsg"` //触发错误信息
	LastSyncTime        time.Time `json:"-"`                   //记录最新一次同步操作时间
	LastWorkTime        time.Time `json:"-"`

	IPAddrHistory      []any `json:"IPAddrHistory"`
	WebhookCallHistroy []any `json:"WebhookCallHistroy"`
	ModifyTime         int64
}

type IPAddrHistoryItem struct {
	IPaddr     string
	RecordTime string
}

type WebhookCallHistroyItem struct {
	CallTime   string
	CallResult string
}

func (d *DDNSTaskState) SetIPAddr(ipaddr string) {
	if d.IpAddr == ipaddr {
		return
	}

	d.IpAddr = ipaddr

	hi := IPAddrHistoryItem{IPaddr: ipaddr, RecordTime: time.Now().Local().Format("2006-01-02 15:04:05")}
	d.IPAddrHistory = append(d.IPAddrHistory, hi)

	if len(d.IPAddrHistory) > 10 {
		d.IPAddrHistory = DeleteAnyListlice(d.IPAddrHistory, 0)
	}
}

func DeleteAnyListlice(a []any, deleteIndex int) []any {
	j := 0
	for i := range a {
		if i != deleteIndex {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func (d *DDNSTaskState) SetDomainUpdateStatus(status string, message string) {
	for i := range d.Domains {
		d.Domains[i].SetDomainUpdateStatus(status, message)
	}
}

func (d *DDNSTaskState) SetWebhookResult(result bool, errMsg string) {
	d.WebhookCallResult = result
	d.WebhookCallErrorMsg = errMsg
	d.WebhookCallTime = time.Now().Format("2006-01-02 15:04:05")

	cr := "成功"
	if !result {
		cr = "出错"
	}

	hi := WebhookCallHistroyItem{CallResult: cr, CallTime: time.Now().Local().Format("2006-01-02 15:04:05")}
	d.WebhookCallHistroy = append(d.WebhookCallHistroy, hi)
	if len(d.WebhookCallHistroy) > 10 {
		d.WebhookCallHistroy = DeleteAnyListlice(d.WebhookCallHistroy, 0)
	}
}

func (d *DDNSTaskState) Init(domains []string, mt int64) {
	d.Domains, d.domainsRawStrList = checkParseDomains(domains)
	d.ModifyTime = mt
}

// Check 检测IP是否有改变
func (d *DDNSTaskState) IPChanged(newAddr string) bool {
	if newAddr == "" {
		return true
	}
	// 地址改变
	if d.IpAddr != newAddr {
		//log.Printf("公网地址改变:[%s]===>[%s]", d.DomainsInfo.IpAddr, newAddr)
		//domains.IpAddr = newAddr
		return true
	}

	return false
}
