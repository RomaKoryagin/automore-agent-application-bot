create table questions (
    id integer primary key,
    chat_id integer not null,
    telegram_id varchar(255) not null,
    descript text null,
    phone text varchar(255) null,
    sended_to_telegram bool default false,
    sended_to_bitfix bool default false
);