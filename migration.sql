CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
create table restaurants
(
    id             uuid primary key DEFAULT uuid_generate_v1(),
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    name           varchar,
    phone          varchar,
    slug           varchar,
    street         varchar,
    district       varchar,
    city           varchar,
    full_address   varchar,
    lat            double precision,
    long           double precision,
    is_opening     boolean,
    is_opening_24h boolean,
    closed         boolean,
    avatar         varchar,
    constraint restaurants_unique_key
        unique (name, full_address)
);