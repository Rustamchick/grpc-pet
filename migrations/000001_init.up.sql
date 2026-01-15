CREATE TABLE grpc_users (
    id              bigserial primary key not null, 
    email           varchar(255) unique not null,
    password_hash   varchar(255) not null,
    is_admin        boolean not null default true
);

CREATE TABLE apps (
    id              bigserial primary key not null, 
    name            varchar(255) unique not null,
    token          varchar(255) not null
);

INSERT INTO apps (id, name, token) VALUES (1, 'rest-todo', 'rest-token');