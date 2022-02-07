-- URL
create or replace function insert_visits_url_if_not_exist(input_url text)
    returns int
    language plpgsql
as
$$
declare
    ua_id integer;
begin
    insert into visits_url (url) values (input_url) returning id into ua_id;
    return ua_id;
exception
    when unique_violation then
        select id from visits_url where url = input_url into ua_id;
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
    insert into visits_ua (ua) values (input_ua) returning id into ua_id;
    return ua_id;
exception
    when unique_violation then
        select id from visits_ua where ua = input_ua into ua_id;
        return ua_id;

end;
$$;

-- IP
create or replace function insert_visits_ip_if_not_exist(input_address text,
                                                         input_bot boolean,
                                                         input_data_center boolean,
                                                         input_tor boolean,
                                                         input_proxy boolean,
                                                         input_vpn boolean,
                                                         input_country text,
                                                         input_domain_count text,
                                                         input_domain_list boolean
)
    returns int
    language plpgsql
as
$$
declare
    ip_id integer;
begin
    insert into visits_ip (address, bot, data_center, tor, proxy, vpn, country, domain_count, domain_list)
    values (input_address::inet, input_bot, input_data_center, input_tor, input_proxy, input_vpn, input_country,
            input_domain_count, input_domain_list)
    returning id into ip_id;
    return ip_id;
exception
    when unique_violation then
        select id from visits_ip where address = input_address::inet into ip_id;
        return ip_id;

end;
$$;

-- API_KEY

create or replace function insert_visits_api_keys_if_not_exist(input_api_key text)
    returns int
    language plpgsql
as
$$
declare
    api_key_id integer;
begin
    insert into visits_api_keys (key) values (input_api_key) returning id into api_key_id;
    return api_key_id;
exception
    when unique_violation then
        select id from visits_api_keys where key = input_api_key into api_key_id;
        return api_key_id;

end;
$$;

-- ACCOUNT
create or replace function insert_visits_account_if_not_exist(input_user_id text,input_api_key integer)
    returns int
    language plpgsql
as
$$
declare
    acc_id integer;
begin
    insert into visits_accounts (user_id,api_key) values (input_user_id,input_api_key) returning id into acc_id;
    return acc_id;
exception
    when unique_violation then
        select id from visits_accounts where user_id = input_user_id into acc_id;
        return acc_id;

end;
$$;