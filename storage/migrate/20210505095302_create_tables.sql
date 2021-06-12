-- +goose Up

create database werewolves;

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
    relations     integer[]               not null default '{}',
    avatar        varchar(68)             not null default 'guest',
    banned_until  timestamp               not null default now() - interval '1' day,
    wins          INTEGER                 not null default 0,
    losses        INTEGER                 not null default 0
);

--friends

create table relationship
(
    created_at    timestamp               not null default now(),
    updated_at    timestamp               not null default now(),
    id      bigserial    not null,
    user_id bigint                not null
);


-- honors

create table honors
(
    id         bigserial primary key   not null,
    created_at    timestamp               not null default now(),
    updated_at    timestamp               not null default now(),
    honored_id bigint      not null,
    reason     varchar(50) not null,
    sender_id  bigint      not null
);

-- reports

create table reports
(
    id          bigserial   not null,
    created_at    timestamp               not null default now(),
    updated_at    timestamp               not null default now(),
    reported_id bigint      not null,
    reason      varchar(50) not null,
    sender_id   bigint      not null
);

-- inventory

create table items
(
    created_at    timestamp               not null default now(),
    updated_at    timestamp               not null default now(),
    id      bigserial   not null,
    name    varchar(50) not null,
    user_id bigint      not null
);

-- triggers
create or replace function trigger_set_timestamp() returns trigger as $$ begin new.updated_at = now();
return new;
end;
$$ language 'plpgsql';


-- +goose Down

