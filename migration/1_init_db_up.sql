
CREATE TABLE if not exists currencies(
    id bigserial primary key,
    name text not null,
    price double precision,
    percent_change_1h double precision,
    percent_change_24h double precision,
    percent_change_7d double precision,
    percent_change_30d double precision,
    last_updated timestamp with time zone
)