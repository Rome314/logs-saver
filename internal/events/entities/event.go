package eventEntities

import (
	"encoding/json"
	"net"
	"time"
)

type RawEvent struct {
	UserId      string  `json:"id" form:"id"`
	Ip          string  `json:"i" form:"i"`
	ApiKey      string  `json:"k" form:"k"`
	Url         string  `json:"u" form:"u"`
	UserAgent   string  `json:"a" form:"a"`
	RequestTime float64 `json:"t" form:"t"`
}

type Event struct {
	Url         string    `json:"url"`
	UserId      string    `json:"user_id"`
	Ip          net.IP    `json:"ip"`
	ApiKey      string    `json:"api_key"`
	Device      Device    `json:"device"`
	IpInfo      *IpInfo   `json:"ip_info"`
	RequestTime time.Time `json:"request_time"`
}
type IpInfo struct {
	Id          int32    `json:"id"`
	Bot         bool     `json:"bot"`
	Datacenter  bool     `json:"datacenter"`
	Tor         bool     `json:"tor"`
	Proxy       bool     `json:"proxy"`
	Vpn         bool     `json:"vpn"`
	Country     string   `json:"country"`
	DomainCount bool     `json:"domaincount,omitempty"`
	DomainList  []string `json:"domainlist,omitempty"`
}

func (i *IpInfo) UnmarshalJSON(bytes []byte) error {
	tmp := struct {
		Id          int32           `json:"id"`
		Bot         bool            `json:"bot"`
		Datacenter  bool            `json:"datacenter"`
		Tor         bool            `json:"tor"`
		Proxy       bool            `json:"proxy"`
		Vpn         bool            `json:"vpn"`
		Country     string          `json:"country"`
		DomainCount json.RawMessage `json:"domaincount,omitempty"`
		DomainList  json.RawMessage `json:"domainlist,omitempty"`
	}{}

	if err := json.Unmarshal(bytes, &tmp); err != nil {
		return err
	}

	i.Id = tmp.Id
	i.Bot = tmp.Bot
	i.Datacenter = tmp.Datacenter
	i.Tor = tmp.Tor
	i.Proxy = tmp.Proxy
	i.Vpn = tmp.Vpn
	i.Country = tmp.Country
	i.DomainList = []string{}

	var cnt bool
	if err := json.Unmarshal(tmp.DomainCount, &cnt); err != nil {
		return nil
	}
	i.DomainCount = cnt

	domainsTmp := [][]string{}
	if err := json.Unmarshal(tmp.DomainList, &domainsTmp); err != nil {
		return err
	}
	domains := []string{}
	for _, d := range domainsTmp {
		if len(d) == 0 {
			continue
		}
		domains = append(domains, d[0])
	}

	i.DomainList = domains
	return nil
}

type Device struct {
	// 0 - stands for desktop, 1 - mobile, 2 - tablet
	Type      int    `json:"type"`
	UserAgent string `json:"user_agent"`
}
