CREATE TABLE IF NOT EXISTS address(
    id serial primary key,
    customer_id int references customer(id),
    street varchar(30)
);
