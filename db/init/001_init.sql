-- Для работы тестов оставим базу по умолчанию postgres
create user api password 'pwd';

create database mart
       owner 'api'
       encoding 'UTF8'
       lc_collate = 'en_US.utf8'
       lc_ctype = 'en_US.utf8';

