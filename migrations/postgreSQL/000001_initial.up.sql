create table visits_ip
(
    id           serial
        constraint visits_ip_pk
            primary key,
    address      inet    not null,
    bot          boolean not null,
    data_center  boolean not null,
    tor          boolean not null,
    proxy        boolean not null,
    vpn          boolean not null,
    country      text    not null,
    domain_count text,
    domain_list  boolean not null
);

comment on table visits_ip is 'Conaints ip information';

alter table visits_ip
    owner to u6lhrq6rcqa0bd;

grant select on sequence visits_ip_id_seq to postgrest_test;

create unique index visits_ip_address_uindex
    on visits_ip (address);

create unique index visits_ip_id_uindex
    on visits_ip (id);

grant select on visits_ip to postgrest_test;

create table visits_url
(
    id  serial
        constraint visits_url_pk
            primary key,
    url text not null
);

comment on table visits_url is 'Visited urls ';

alter table visits_url
    owner to u6lhrq6rcqa0bd;

grant select on sequence visits_url_id_seq to postgrest_test;

create unique index visits_url_url_uindex
    on visits_url (url);

create unique index visits_url_id_uindex
    on visits_url (id);

grant select on visits_url to postgrest_test;

create table visits_ua
(
    id serial
        constraint visits_ua_pk
            primary key,
    ua text not null
);

comment on table visits_ua is 'Conains user agents';

alter table visits_ua
    owner to u6lhrq6rcqa0bd;

grant select on sequence visits_ua_id_seq to postgrest_test;

create unique index visits_ua_ua_uindex
    on visits_ua (ua);

create unique index visits_ua_id_uindex
    on visits_ua (id);

grant select on visits_ua to postgrest_test;

create table visits_api_keys
(
    id       smallserial
        constraint visits_api_keys_pk
            primary key,
    key      text                  not null,
    "quote " integer default 50000 not null
);

comment on table visits_api_keys is 'List of api kets with quotas';

alter table visits_api_keys
    owner to u6lhrq6rcqa0bd;

grant select on sequence visits_api_keys_id_seq to postgrest_test;

create unique index visits_api_keys_id_uindex
    on visits_api_keys (id);

create unique index visits_api_keys_key_uindex
    on visits_api_keys (key);

grant select on visits_api_keys to postgrest_test;

create table visits_accounts
(
    id           serial
        constraint visits_accounts_pk
            primary key,
    user_id      text                                not null,
    ips          integer[] default '{}'::integer[]   not null,
    countries    integer[] default '{}'::integer[]   not null,
    total_visits integer   default 0                 not null,
    api_key      integer                             not null
        constraint visits_accounts_api_key_fk
            references visits_api_keys,
    last_ip      inet,
    created      timestamp default CURRENT_TIMESTAMP not null,
    last_updated timestamp default CURRENT_TIMESTAMP not null
);

alter table visits_accounts
    owner to u6lhrq6rcqa0bd;

grant select on sequence visits_accounts_id_seq to postgrest_test;

create unique index visits_accounts_id_uindex
    on visits_accounts (id);

create unique index visits_accounts_user_id_uindex
    on visits_accounts (user_id);

grant select on visits_accounts to postgrest_test;

create table visits
(
    id      bigserial
        constraint visits_pk
            primary key,
    api_key integer   not null
        constraint visits_api_key__fk
            references visits_api_keys,
    account integer   not null
        constraint visits_account_fk
            references visits_accounts,
    ip      integer   not null
        constraint visits_ip_fk
            references visits_ip,
    url     integer   not null
        constraint visits_url_fk
            references visits_url,
    ua      integer   not null
        constraint visits_ua_fk
            references visits_ua,
    time    timestamp not null
);

alter table visits
    owner to u6lhrq6rcqa0bd;

grant select on sequence visits_id_seq to postgrest_test;

create unique index visits_id_uindex
    on visits (id);

grant select on visits to postgrest_test;

