drop table if exists balance;
create table balance (
    user_id int references usr(id),
    available float8 not null default 0,
    withdrawn float8 not null default 0,
    constraint current_positive check (available >= 0),
    constraint withdraw_not_negative check(withdrawn >= 0)
);

create unique index ixc_balance_user_id on balance (user_id);
cluster balance using ixc_balance_user_id;

insert into balance (user_id, available, withdrawn)
select id, 0, 0
from usr
where id not in (select user_id from balance);


