create table if not exists logs
(
    id           serial
        constraint logs_pk
            primary key,
    user_id      int not null,
    api_key      varchar(64),
    url          text,
    user_agent   text,
    request_time timestamp,
    ip           varchar(15)
);

create index if not exists logs_api_key_index
    on logs (api_key);

create unique index if not exists logs_id_uindex
    on logs (id);

create index if not exists logs_request_time_index
    on logs (request_time desc);

create index if not exists logs_user_id_index
    on logs (user_id);

