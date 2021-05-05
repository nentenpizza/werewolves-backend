-- +goose Up

-- users

create table users
(
    created_at    timestamp               not null default now(),
    updated_at    timestamp               not null default now(),
    xp            bigint                  not null default 0,
    id            bigserial               not null,
    email         varchar(255)            not null,
    login         varchar(25)             not null,
    username      varchar(13) primary key not null,
    password_hash varchar(60)             not null,
    avatar        varchar(68)             not null default 'guest',
    banned_until  timestamp               not null default now() - interval '1' day,
    wins          INTEGER                 not null default 0,
    losses        INTEGER                 not null default 0
);

--friends

create table friends
(
    id         bigserial not null,
    created_at timestamp not null default now(),
    sender_id  bigint    not null,
    target_id  bigint    not null,
    active     bool      not null
);

-- honors

create table honors
(
    id         bigserial   not null,
    created_at timestamp   not null default now(),
    honored_id bigint      not null,
    reason     varchar(50) not null,
    sender_id  bigint      not null
);

-- reports

create table reports
(
    id          bigserial   not null,
    created_at  timestamp   not null default now(),
    reported_id bigint      not null,
    reason      varchar(50) not null,
    sender_id   bigint      not null
);

-- inventory

create table inventory
(
    id      bigserial   not null,
    item    varchar(50) not null,
    user_id bigint      not null
);

-- triggers


-- +goose Down

