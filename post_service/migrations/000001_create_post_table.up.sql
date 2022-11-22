CREATE TABLE IF NOT EXISTS post(
    id serial primary key,
    name varchar(30),
    description text,
    customer_id int,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
)