CREATE TABLE IF NOT EXISTS media(
    id serial primary key,
    name varchar(30),
    post_id int references post(id)
)