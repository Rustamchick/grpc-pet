CREATE TABLE users (
    id              bigserial primary key not null, 
    email           varchar(255) unique not null,
    password_hash   varchar(255) not null
);

