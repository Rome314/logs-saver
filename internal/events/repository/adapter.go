package eventsRepository

import (
	"database/sql"
	"sync"

	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
)

type eventSql struct {
	Url         sql.NullString `db:"url"`
	UserId      sql.NullInt64  `db:"user_id"`
	Ip          sql.NullString `db:"ip"`
	ApiKey      sql.NullString `db:"api_key"`
	UserAgent   sql.NullString `db:"user_agent"`
	RequestTime sql.NullTime   `db:"request_time"`
}

func eventToSql(input *eventEntities.Event) eventSql {
	return eventSql{
		Url:         sql.NullString{input.Url, true},
		UserId:      sql.NullInt64{input.UserId, true},
		Ip:          sql.NullString{input.Ip.String(), true},
		ApiKey:      sql.NullString{input.ApiKey, true},
		UserAgent:   sql.NullString{input.UserAgent, true},
		RequestTime: sql.NullTime{input.RequestTime, true},
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
