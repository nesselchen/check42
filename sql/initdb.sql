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

create table if not exists `todo`(
    `id` int not null auto_increment,
    `owner` int not null,
    `text` varchar(140),
    `done` boolean default 0,
    `created` datetime default current_timestamp,
    primary key (`id`),
    foreign key (`owner`) references user(`id`)
);

insert into user (`name`, `email`, `password_hash`) values
    ("admin", "admin@adm.in", "$2a$10$tYtiAnVF2EJ6WPT894/YaO.VoQ08sknhVSa2jT0Sac1bvK2AgWeN.");
    
insert into todo (`owner`, `text`) values
    (1, "Input validation"),
    (1, "User authentication"),
    (1, "Frontend client"),
    (1, "Testing");