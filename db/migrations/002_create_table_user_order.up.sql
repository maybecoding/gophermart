drop type if exists order_status;
create type order_status as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

drop table if exists user_order;
create table user_order (
                            user_id int references usr(id),
                            number varchar(255) unique not null,
                            status order_status not null,
                            accrual int default 0,
                            created_at date not null default now()
);