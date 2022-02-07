package eventEntities

import (
	"net"
	"time"
)

// ‘i’ => $ip
// ‘k’ => $apikey
// ‘u’ => urlencode($link20_idkb)
// ‘a’ => $ua
// ’t’ => $request_time
// ‘id’ => $user->id

type RawEvent struct {
	UserId      string `json:"id"`
	Ip          string `json:"i"`
	ApiKey      string `json:"k"`
	Url         string `json:"u"`
	UserAgent   string `json:"a"`
	RequestTime int64  `json:"t"`
}

type Event struct {
	Url         string    `json:"url"`
	UserId      string    `json:"user_id"`
	Ip          net.IP    `json:"ip"`
	ApiKey      string    `json:"api_key"`
	UserAgent   string    `json:"user_agent"`
	IpInfo      *IpInfo   `json:"ip_info"`
	RequestTime time.Time `json:"request_time"`
}
type IpInfo struct {
	Id          int32  `json:"id"`
	Bot         bool   `json:"bot"`
	Datacenter  bool   `json:"datacenter"`
	Tor         bool   `json:"tor"`
	Proxy       bool   `json:"proxy"`
	Vpn         bool   `json:"vpn"`
	Country     string `json:"country"`
	DomainCount string `json:"domaincount"`
	DomainList  string `json:"domain_list"`
}
