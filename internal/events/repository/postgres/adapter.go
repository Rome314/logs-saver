package events_postgres_repository

import (
	"fmt"
	"strings"
	"sync"

	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
)

func GetQueryValueMany(events []*eventEntities.Event) []string {
	count := len(events)

	resp := make([]string, count)

	wg := &sync.WaitGroup{}
	wg.Add(count)

	for i, e := range events {
		go func(index int, event *eventEntities.Event) {
			defer wg.Done()
			resp[index] = GetQueryValue(event)
		}(i, e)
	}
	wg.Wait()
	return resp
}

func sanitize(event *eventEntities.Event) {
	event.Url = strings.ReplaceAll(event.Url, "'", `"`)
	event.ApiKey = strings.ReplaceAll(event.ApiKey, "'", `"`)
	event.UserId = strings.ReplaceAll(event.UserId, "'", `"`)
	event.Device.UserAgent = strings.ReplaceAll(event.Device.UserAgent, "'", `"`)
}

func GetQueryValue(event *eventEntities.Event) string {
	sanitize(event)
	ipInfo := event.IpInfo
	return fmt.Sprintf(
		`(
		        insert_visits_api_keys_if_not_exist('%s'),
        		insert_visits_account_if_not_exist(
                		'%s',
                		insert_visits_api_keys_if_not_exist('%s')),
		        case 
		            when %d != 0 
					then %d 
		            else 
						insert_visits_ip_if_not_exist(
								'%s',
								%t,
								%t,
								%t,
								%t,
								%t, 
						    	'%s',
								{%s})
				end ,
				insert_visits_url_if_not_exist('%s'),
				insert_visits_device_if_not_exist(%s,%s,%s,%d,TO_TIMESTAMP('%s','YYYY-MM-DD HH24:MI:SS.US')),
				TO_TIMESTAMP('%s','YYYY-MM-DD HH24:MI:SS.US'))`,
		event.ApiKey,
		event.UserId,
		event.ApiKey,
		event.IpInfo.Id,
		event.IpInfo.Id,
		event.Ip.String(),
		ipInfo.Bot,
		ipInfo.Datacenter,
		ipInfo.Tor,
		ipInfo.Proxy,
		ipInfo.Vpn,
		ipInfo.Country,
		strings.Join(event.IpInfo.DomainList, ","),
		event.Url,
		event.UserId,
		event.Device.UserAgent,
		event.ApiKey,
		event.Device.Type,
		event.RequestTime.Format("2006-01-02 15:04:05.000000"),
		event.RequestTime.Format("2006-01-02 15:04:05.000000"),
	)

}
