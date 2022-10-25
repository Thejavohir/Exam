CREATE TABLE IF NOT EXISTS customer(
    id serial primary key,
    first_name varchar(30),
    last_name varchar(30),
    bio text,
    email varchar(30),
    phone_number varchar(15),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
