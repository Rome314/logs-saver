package events

import (
	"net"
	"net/url"
	"strings"
	"time"

	"emperror.dev/errors"
	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
)

func RawToEvent(input eventEntities.RawEvent) (event *eventEntities.Event, err error) {
	if input.UserId == "" {
		err = errors.New("user id not provided")
		return
	}

	if input.ApiKey == "" {
		err = errors.New("api key not provided")
		return
	}

	if input.UserAgent == "" {
		err = errors.New("user agent not provided")
		return
	}

	decodedUrl, err := url.QueryUnescape(input.Url)
	if err != nil {
		err = errors.WithMessage(err, "decoding input url")
		return
	}

	u, err := url.ParseRequestURI(decodedUrl)
	if err != nil {
		err = errors.WithMessage(err, "validating input url")
		return
	}

	splitedIp := strings.Split(input.Ip, "/")
	rawIp := splitedIp[0]
	if len(splitedIp) == 2 {
		rawIp = splitedIp[1]
	}

	ip := net.ParseIP(rawIp)
	if ip == nil {
		err = errors.Errorf("invalid ip provided: %s", rawIp)
		return
	}

	requestTime := time.UnixMicro(input.RequestTime)
	if requestTime.IsZero() {
		err = errors.New("invalid request time provided")
		return
	}

	event = &eventEntities.Event{
		Url:         u.String(),
		UserId:      input.UserId,
		Ip:          ip,
		ApiKey:      input.ApiKey,
		UserAgent:   input.UserAgent,
		RequestTime: requestTime,
	}
	return event, nil
}
