-- URL
create or replace function insert_visits_url_if_not_exist(input_url text)
    returns int
    language plpgsql
as
$$
declare
    ua_id integer;
begin
    select id from visits_url where url = input_url into ua_id;
    if ua_id IS NULL then
        insert into visits_url (url) values (input_url) returning id into ua_id;
    end if;
    return ua_id;
end;
$$;
-- USER_AGENT
create or replace function insert_visits_ua_if_not_exist(input_ua text)
    returns int
    language plpgsql
as
$$
declare
    ua_id integer;
begin
    select id from visits_ua where ua = input_ua into ua_id;
    if ua_id IS NULL then
        insert into visits_ua (ua) values (input_ua) returning id into ua_id;
    end if;
    return ua_id;

end;
$$;

-- IP
]

-- API_KEY

create or replace function insert_visits_api_keys_if_not_exist(input_api_key text)
    returns int
    language plpgsql
as
$$
declare
    api_key_id integer;
begin
    select id from visits_api_keys where key = input_api_key into api_key_id;
    if api_key_id IS NULL then
        insert into visits_api_keys (key) values (input_api_key) returning id into api_key_id;
    end if;
    return api_key_id;


end;
$$;

-- ACCOUNT
create or replace function insert_visits_account_if_not_exist(input_user_id text, input_api_key integer)
    returns int
    language plpgsql
as
$$
declare
    acc_id integer;
begin
    select id from visits_accounts where user_id = input_user_id into acc_id;
    if acc_id IS NULL then
        insert into visits_accounts (user_id, api_key) values (input_user_id, input_api_key) returning id into acc_id;
    end if;
    return acc_id;

end;
$$;

-- DEVICE
create or replace function insert_visits_device_if_not_exist(input_user text, input_ua text, input_key text,
                                                             input_type smallint, input_time timestamp) returns integer
    language plpgsql
as
$$
declare
    device_id  integer;
    account_id integer;
    ua_id      integer;
    key_id     integer;

begin

    select * from insert_visits_api_keys_if_not_exist(input_key) into key_id;
    select * from insert_visits_ua_if_not_exist(input_ua) into ua_id;
    select * from insert_visits_account_if_not_exist(input_user, key_id) into account_id;

    select id from visits_devices where ua = ua_id and key = key_id and type = input_type into device_id;
    if device_id IS NULL then
        insert into visits_devices(account_id, type, key, ua, created)
        values (account_id, input_type, key_id, ua_id, input_time)
        returning id into device_id;
    end if;
    return device_id;

end;
$$;


