CREATE TABLE platforms (
  id          varchar(100) not null unique,
  name        varchar(100) not null unique,
  type        varchar(50) not null,
  description varchar(500)
);
