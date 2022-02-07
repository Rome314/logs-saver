package events_postgres_repository

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

func (r *repo) Status() error {
	return r.client.Ping()
}

func NewPostgres(logger *logging.Entry, client *sqlx.DB) eventEntities.Repository {
	return &repo{logger: logger, client: client}
}

func (r *repo) Store(event *eventEntities.Event) (err error) {
	query := `
        insert into visits(api_key, account, ip, url, ua, time)
		VALUES (
		        insert_visits_api_keys_if_not_exist(:api_key),
        		insert_visits_account_if_not_exist(
                		:user_id,
                		insert_visits_api_keys_if_not_exist(:api_key)),
		        case 
		            when :ip_info.id != 0 
					then :ip_info.id 
		            else 
						insert_visits_ip_if_not_exist(
								:ip,
								:ip_info.bot,
								:ip_info.data_center,
								:ip_info.tor,
								:ip_info.proxy,
								:ip_info.vpn, 
						    	:ip_info.country,
								:ip_info.domain_count,
						    	:ip_info.domain_list)
				end ,
				insert_visits_url_if_not_exist(:url),
				insert_visits_ua_if_not_exist(:user_agent),
				:request_time)`

	tmp := eventToSql(event)

	_, err = r.client.NamedExec(query, tmp)
	if err != nil {
		err = errors.WithMessage(err, "executing query")
		return
	}

	return nil
}
func (r *repo) StoreMany(events ...*eventEntities.Event) (inserted int64, err error) {
	// logger := r.logger.WithMethod("StoreMany")

	query := `
        insert into visits(api_key, account, ip, url, ua, time)
		VALUES (
		        insert_visits_api_keys_if_not_exist(:api_key),
        		insert_visits_account_if_not_exist(
                		:user_id,
                		insert_visits_api_keys_if_not_exist(:api_key)),
		        case 
		            when :ip_info.id != 0 
					then :ip_info.id 
		            else 
						insert_visits_ip_if_not_exist(
								:ip,
								:ip_info.bot,
								:ip_info.data_center,
								:ip_info.tor,
								:ip_info.proxy,
								:ip_info.vpn, 
						    	:ip_info.country,
								:ip_info.domain_count,
						    	:ip_info.domain_list)
				end ,
				insert_visits_url_if_not_exist(:url),
				insert_visits_ua_if_not_exist(:user_agent),
				:request_time)`

	eventsSql := eventToSqlMany(events...)

	res, err := r.client.NamedExec(query, eventsSql)
	if err != nil {
		err = errors.WithMessage(err, "executing query")
		return
	}
	inserted, _ = res.RowsAffected()

	return inserted, nil

}
