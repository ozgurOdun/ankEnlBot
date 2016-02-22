# ankEnlBot

This bot gets token from a file located at the root named token.env
Also uses a database named new.db located under database folder

CREATE TABLE `Task` (
	`Uid`	INTEGER UNIQUE,
	`task`	TEXT,
	`creation_time`	NUMERIC,
	`ownerid`	INTEGER,
	`done`	INTEGER,
	PRIMARY KEY(Uid)
);

CREATE TABLE `User` (
	`Uid`	INTEGER UNIQUE,
	`user_name`	TEXT UNIQUE,
	`first_name`	TEXT,
	`last_name`	TEXT,
	`creation_time`	NUMERIC,
	PRIMARY KEY(Uid)
);
