package eventsRedisRepository

import (
	"context"
	"encoding/json"
	"sync"

	"emperror.dev/errors"
	"github.com/go-redis/redis/v8"
	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
	log "github.com/sirupsen/logrus"
)

type repo struct {
	client *redis.Client
}

func New(client *redis.Client) eventEntities.BufferRepo {
	return &repo{
		client: client,
	}
}

func (r *repo) StoreToErrorStorage(events []*eventEntities.Event) (err error) {

	values := make([]interface{}, len(events))

	wg := &sync.WaitGroup{}
	wg.Add(len(events))

	for index, event := range events {
		go func(i int, e *eventEntities.Event) {
			defer wg.Done()
			bts, _ := json.Marshal(e)
			values[i] = string(bts)
		}(index, event)
	}
	wg.Wait()

	if err = r.client.SAdd(context.TODO(), "uninserteds_buffer", values...).Err(); err != nil {
		return errors.WithMessage(err, "inserting")
	}
	return nil

}

func (r *repo) Count() (count uint64, err error) {
	count, err = r.client.SCard(context.TODO(), "insert_buffer").Uint64()
	if err != nil {
		err = errors.WithMessage(err, "getting count")
		return
	}
	return count, nil
}

func (r *repo) PopAll() (events []*eventEntities.Event, err error) {
	count, err := r.Count()
	if err != nil {
		err = errors.WithMessage(err, "getting count")
		return
	}

	members, err := r.client.SPopN(context.TODO(), "insert_buffer", int64(count)).Result()
	if err != nil {
		err = errors.WithMessage(err, "getting")
	}

	events = make([]*eventEntities.Event, len(members))

	wg := &sync.WaitGroup{}
	wg.Add(len(members))

	for index, str := range members {
		go func(i int, encoded string) {
			defer wg.Done()
			tmp := &eventEntities.Event{}
			if e := json.Unmarshal([]byte(encoded), tmp); e != nil {
				log.Error(e)
				return
			}
			events[i] = tmp

		}(index, str)
	}
	wg.Wait()

	return events, nil

}
func (r *repo) Store(event *eventEntities.Event) (bufferSize uint64, err error) {
	bts, err := json.Marshal(event)
	if err != nil {
		err = errors.WithMessage(err, "marshalling")
		return
	}

	if err = r.client.SAdd(context.TODO(), "insert_buffer", string(bts)).Err(); err != nil {
		err = errors.WithMessage(err, "inserting")
		return
	}

	return r.Count()
}

func (r *repo) Status() error {
	return r.client.Ping(context.TODO()).Err()
}
