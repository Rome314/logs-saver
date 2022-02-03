package eventsRepository

import (
	"emperror.dev/errors"
	"github.com/jmoiron/sqlx"
	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
	"github.com/rome314/idkb-events/pkg/logging"
)

type repo struct {
	logger *logging.Entry
	client *sqlx.DB
}

func NewPostgres(logger *logging.Entry, client *sqlx.DB) eventEntities.Repository {
	return &repo{logger: logger, client: client}
}

func (r *repo) StoreMany(events ...*eventEntities.Event) (inserted int64, err error) {
	// logger := r.logger.WithMethod("StoreMany")

	query := `insert into logs ( user_id, api_key, url, user_agent, request_time, ip)
				values (:user_id, :api_key, :url, :user_agent, :request_time, :ip)
				on conflict do nothing;`

	eventsSql := eventToSqlMany(events...)

	res, err := r.client.NamedExec(query, eventsSql)
	if err != nil {
		err = errors.WithMessage(err, "executing query")
		return
	}
	inserted, _ = res.RowsAffected()

	//
	// tx, err := r.client.Beginx()
	// if err != nil {
	// 	err = errors.WithMessage(err, "starting transaction")
	// 	return
	// }
	// defer tx.Rollback()
	//
	// stmt, err := tx.PrepareNamed(query)
	// if err != nil {
	// 	err = errors.WithMessage(err, "preparing statement")
	// 	return
	// }
	//
	// for _, event := range eventsSql {
	// 	_, e := stmt.Exec(event)
	// 	if e != nil {
	// 		logger.WithPlace("exec_stmt").Error(e)
	// 		continue
	// 	}
	// 	inserted++
	//
	// }
	// err = stmt.Close()
	// if err != nil {
	// 	err = errors.WithMessage(err, "close stmt")
	// 	return
	// }
	//
	// err = tx.Commit()
	// if err != nil {
	// 	err = errors.WithMessage(err, "committing")
	// 	return
	// }

	return inserted, nil

}
