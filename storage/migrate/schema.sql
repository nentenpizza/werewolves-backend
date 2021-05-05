-- +goose Up

-- users

create table users(
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    xp bigint not null default 0,
    id bigserial not null,
    email varchar(255) not null,
    login varchar(25) unique not null,
    username varchar(13) primary key not null,
    password_hash varchar(60) not null,
    avatar varchar(68) not null default 'guest',
    banned_until timestamp not null default now() - interval '1' day,
    wins INTEGER not null default 0,
    losses INTEGER not null default 0,
    rating integer not null default 0
);

-- reports

create table reports(
    reported_id bigint not null,
    reason varchar(50) not null,
    sender_id bigint not null,
    note varchar(300) not null default ''
);

-- honors

create table honors(
  honored_id bigint not null,
  reason varchar(200) not null,
  type varchar(10) not null,
  sender_id bigint not null
);

create table inventory(
    item varchar(64) not null,
    user_id bigint not null,
    count integer not null
);

-- triggers
create or replace function trigger_set_timestamp() returns trigger as $$ begin new.updated_at = now();
return new;
end;
$$ language 'plpgsql';