drop table if exists "order";
drop type if exists order_status;

create type order_status as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');


create table "order" (
                            user_id int references usr(id),
                            order_nr varchar(255) unique not null,
                            status order_status not null,
                            accrual float8 default 0,
                            accrual_at timestamptz null,
                            created_at date not null default now()
);