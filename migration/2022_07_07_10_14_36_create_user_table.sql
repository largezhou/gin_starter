-- __UP__ 请勿删除该行
create table user
(
    id          bigint unsigned auto_increment,
    uuid        varchar(36)  not null,
    username    varchar(50)  not null,
    password    varchar(100) not null,
    email       varchar(100) null,
    create_time datetime     not null,
    update_time datetime     not null,
    constraint user_pk
        primary key (id)
);

create unique index user_username_uindex
    on user (username);

create unique index user_uuid_uindex
    on user (uuid);

-- __DOWN__ 请勿删除该行
drop table user;
