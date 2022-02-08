package eventsRedisRepository

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
)

func hash(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func getKeyValue(event *eventEntities.Event) (key string, value string) {
	key = fmt.Sprintf("%s_%s_%d", event.ApiKey, event.UserId, event.RequestTime.UnixMilli())
	key = hash(key)
	bts, _ := json.Marshal(event)
	value = string(bts)
	return
}
