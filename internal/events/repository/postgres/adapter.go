package events_postgres_repository

import (
	"database/sql"
	"sync"

	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
)

type ipInfoSql struct {
	Id          sql.NullInt32  `db:"id"`
	Bot         sql.NullBool   `db:"bot"`
	Datacenter  sql.NullBool   `db:"data_center"`
	Tor         sql.NullBool   `db:"tor"`
	Proxy       sql.NullBool   `db:"proxy"`
	Vpn         sql.NullBool   `db:"vpn"`
	Country     sql.NullString `db:"country"`
	DomainCount sql.NullString `db:"domain_count"`
	DomainList  sql.NullBool   `db:"domain_list"`
	Address     sql.NullString `db:"address"`
}

func (i ipInfoSql) ToIpInfo() *eventEntities.IpInfo {
	return &eventEntities.IpInfo{
		Id:          i.Id.Int32,
		Bot:         i.Bot.Bool,
		Datacenter:  i.Datacenter.Bool,
		Tor:         i.Tor.Bool,
		Proxy:       i.Proxy.Bool,
		Vpn:         i.Vpn.Bool,
		Country:     i.Country.String,
		DomainCount: i.DomainCount.String,
		DomainList:  "",
	}
}

type eventSql struct {
	Url         sql.NullString `db:"url"`
	UserId      sql.NullString `db:"user_id"`
	Ip          sql.NullString `db:"ip"`
	ApiKey      sql.NullString `db:"api_key"`
	UserAgent   sql.NullString `db:"user_agent"`
	IpInfo      ipInfoSql      `db:"ip_info"`
	RequestTime sql.NullTime   `db:"request_time"`
}

func ipInfoToSql(input *eventEntities.IpInfo) ipInfoSql {
	return ipInfoSql{
		Id:          sql.NullInt32{input.Id, true},
		Bot:         sql.NullBool{input.Bot, true},
		Datacenter:  sql.NullBool{input.Datacenter, true},
		Tor:         sql.NullBool{input.Tor, true},
		Proxy:       sql.NullBool{input.Proxy, true},
		Vpn:         sql.NullBool{input.Vpn, true},
		Country:     sql.NullString{input.Country, true},
		DomainCount: sql.NullString{input.DomainCount, true},
		DomainList:  sql.NullBool{false, true},
	}
}

func eventToSql(input *eventEntities.Event) eventSql {
	return eventSql{
		Url:         sql.NullString{input.Url, true},
		UserId:      sql.NullString{input.UserId, true},
		Ip:          sql.NullString{input.Ip.String(), true},
		ApiKey:      sql.NullString{input.ApiKey, true},
		UserAgent:   sql.NullString{input.UserAgent, true},
		RequestTime: sql.NullTime{input.RequestTime, true},
		IpInfo:      ipInfoToSql(input.IpInfo),
	}
}

func eventToSqlMany(events ...*eventEntities.Event) []eventSql {
	resp := make([]eventSql, len(events))

	wg := &sync.WaitGroup{}
	wg.Add(len(events))

	for index, event := range events {
		go func(i int, e *eventEntities.Event) {
			defer wg.Done()
			resp[i] = eventToSql(e)
		}(index, event)
	}
	wg.Wait()
	return resp
}
