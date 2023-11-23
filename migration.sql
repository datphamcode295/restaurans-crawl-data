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
    is_opening_24h boolean,
    avatar         varchar,
    external_id    varchar unique not null
);