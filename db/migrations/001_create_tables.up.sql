drop table if exists usr;
create table usr (
    user_id int primary key generated always as identity,
    login varchar(255) unique,
    hash varchar(255)
);