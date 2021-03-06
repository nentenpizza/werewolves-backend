create database werewolves;

create table users(
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    xp bigint not null default 0,
    id bigserial not null,
    email varchar(255) not null,
    username varchar(15) primary key not null default '',
    password_hash varchar(60) not null,
    avatar varchar(68) not null default 'guest',
    banned_until timestamp not null default now() - interval '1' day,
    wins INTEGER not null default 0,
    losses INTEGER not null default 0
);
-- triggers
create or replace function trigger_set_timestamp() returns trigger as $$ begin new.updated_at = now();
return new;
end;
$$ language 'plpgsql';