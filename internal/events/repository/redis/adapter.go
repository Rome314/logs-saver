package eventsRedisRepository

import (
	"encoding/json"
	"fmt"

	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
)

func getKeyValue(event *eventEntities.Event) (key string, value string) {
	key = fmt.Sprintf("%s_%s_%d", event.ApiKey, event.UserId, event.RequestTime.UnixMilli())
	bts, _ := json.Marshal(event)
	value = string(bts)
	return
}
