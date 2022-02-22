package events_postgres_repository

import (
	"fmt"
	"strings"

	"emperror.dev/errors"
	"github.com/jmoiron/sqlx"
	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
	"github.com/rome314/idkb-events/pkg/logging"
)

type repo struct {
	logger *logging.Entry
	client *sqlx.DB
}

func (r *repo) Status() error {
	return r.client.Ping()
}

func NewPostgres(logger *logging.Entry, client *sqlx.DB) eventEntities.Repository {
	return &repo{logger: logger, client: client}
}

func (r *repo) Store(event *eventEntities.Event) (err error) {
	query := fmt.Sprintf(`insert into visits(api_key, account, ip, url, ua, time)
		VALUES %s;`, GetQueryValue(event))

	_, err = r.client.Exec(query)
	if err != nil {
		err = errors.WithMessage(err, "executing query")
		return
	}

	return nil
}

func (r *repo) StoreMany(events ...*eventEntities.Event) (inserted int64, err error) {
	// logger := r.logger.WithMethod("StoreMany")
	tx, err := r.client.Beginx()
	if err != nil {
		err = errors.WithMessage(err, "creating tx")
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(PreInsertQueries)
	if err != nil {
		err = errors.WithMessage(err, "running pre insert queries")
		return
	}

	values := GetQueryValueMany(events)
	query := fmt.Sprintf(`insert into visits(api_key, account, ip, url, ua, time)
		VALUES %s;`, strings.Join(values, ","))

	res, err := tx.Exec(query)
	if err != nil {
		err = errors.WithMessage(err, "executing inserts")
		return
	}

	_, err = tx.Exec(PostInsertQueries)
	if err != nil {
		err = errors.WithMessage(err, "running post insert queries")
		return
	}

	if err = tx.Commit(); err != nil {
		err = errors.WithMessage(err, "committing tx")
		return
	}
	inserted, _ = res.RowsAffected()

	return inserted, nil
}
