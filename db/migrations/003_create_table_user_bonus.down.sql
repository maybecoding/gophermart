-- alter table user_order drop column id;
drop table if exists user_bonus_balance;
alter table user_order drop column if exists accrual_at;