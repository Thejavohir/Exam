CREATE TABLE IF NOT EXISTS review(
    id serial primary key,
    name varchar(20),
    review int check,
    description text,
    post_id int,
    created_at TIME NOT NULL DEFAULT NOW(), 
    updated_at TIME NOT NULL DEFAULT NOW(),
    deleted_at TIME
)