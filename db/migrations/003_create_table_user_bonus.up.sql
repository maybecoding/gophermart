-- drop table if exists user_order;
-- create table user_order (
--                             id int primary key generated always as identity,
--                             user_id int references usr(id),
--                             number varchar(255) unique not null,
--                             status order_status not null,
--                             accrual int default 0,
--                             created_at date not null default now()
-- );

-- drop table if exists user_bonus;
--
-- create table user_bonus (
--     user_id int references usr(id),
--     order_id int references user_order(id),
--     amount int not null default 0,
--     processed_at date not null default now()
-- );

-- create index ixc_user_bonus_user_id on user_bonus (user_id);
-- cluster user_bonus using ixc_user_bonus_user_id;

drop table if exists user_bonus_balance;
create table user_bonus_balance (
    user_id int references usr(id),
    available int not null default 0,
    withdrawn int not null default 0,
    constraint current_positive check (available >= 0),
    constraint withdraw_not_negative check(withdrawn >= 0)
);

create unique index ixc_user_bonus_balance_user_id on user_bonus_balance (user_id);
cluster user_bonus_balance using ixc_user_bonus_balance_user_id;

insert into user_bonus_balance (user_id, available, withdrawn)
select id, 0, 0
from usr
where id not in (select user_id from user_bonus_balance);

alter table user_order drop column if exists accrual_at;
alter table user_order add accrual_at timestamptz null;

