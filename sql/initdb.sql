create database if not exists `check42` ;

use `check42`;

create table if not exists `user` (
    `id` int not null auto_increment,
    `name` varchar(140) not null,
    `email` varchar(50) not null,
    `password_hash` varchar(255) not null,
    `created` datetime default current_timestamp,
    primary key (`id`),
    unique (`name`),
    unique (`email`)
);

create table if not exists `todo_category` (
	`id` 	int not null auto_increment,
    `name` varchar(140) default "New category",
    `owner` int not null,
    primary key (`id`),
    foreign key (`owner`) references `user` (`id`)
);

create table if not exists `todo` (
    `id`        int not null auto_increment,
    `owner`     int not null,
    `text`      varchar(140),
    `done`      boolean default 0,
    `created`   datetime default current_timestamp,
    `category`  int null,
    primary key (`id`),
    foreign key (`owner`) references `user` (`id`),
    foreign key (`category`) references `todo_category` (`id`) on delete cascade
);

insert into `user` (`name`, `email`, `password_hash`) values
    ("admin", "admin@adm.in", "$2a$10$vPibycoXtT9WGUAEHrF/LeU.X2GM3UC4/mx8av2o63M5rXtQgDsw2");

insert into `todo_category` (`name`, `owner`) values
	("My tasks", 1),
    ("Urgent", 1);

insert into `todo` (`owner`, `text`, `category`) values
    (1, "Laundry", 1),
    (1, "Dishes", 2);

insert into `todo` (`owner`, `text`) values
    (1, "Test"),
    (1, "Todo");
