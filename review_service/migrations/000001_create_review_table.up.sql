CREATE TABLE IF NOT EXISTS review(
    id serial primary key,
    name varchar(30),
    review int check(review>0 and review<6),
    description text,
    post_id int,
    created_at TIME NOT NULL DEFAULT NOW(), 
    updated_at TIME NOT NULL DEFAULT NOW(),
    deleted_at TIME
)