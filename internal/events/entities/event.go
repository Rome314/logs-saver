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
	Ip          string `json:"i"`
	ApiKey      string `json:"k"`
	Url         string `json:"u"`
	UserAgent   string `json:"a"`
	RequestTime int64  `json:"t"`
	UserId      int64  `json:"id"`
}

type Event struct {
	Url         string    `json:"url"`
	UserId      int64     `json:"user_id"`
	Ip          net.IP    `json:"ip"`
	ApiKey      string    `json:"api_key"`
	UserAgent   string    `json:"user_agent"`
	RequestTime time.Time `json:"request_time"`
}
