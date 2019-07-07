CREATE TABLE IF NOT EXISTS viraagh_user
(
  user_id       SERIAL,
  first_name    TEXT      NULL,
  last_name     TEXT      NULL,
  username      TEXT      NOT NULL,
  password      TEXT      NOT NULL,
  linkedIn_url  TEXT      NOT NULL,
  filename      TEXT      NULL,
  image_link    TEXT      NULL,
  insert_time   TIMESTAMP NOT NULL,
  PRIMARY KEY   (user_id),
  UNIQUE        (linkedIn_url)
);

CREATE TABLE IF NOT EXISTS school
(
  school_id       SERIAL,
  school_name     TEXT      NOT NULL,
  degree          TEXT      NULL,
  field_of_study  TEXT      NULL,
  insert_time     TIMESTAMP NOT NULL,
  PRIMARY KEY     (school_id)
);

CREATE TABLE IF NOT EXISTS company
(
  company_id    SERIAL,
  company_name  TEXT      NOT NULL,
  location      TEXT      NULL,
  insert_time   TIMESTAMP NOT NULL,
  PRIMARY KEY   (company_id)
);

CREATE TABLE IF NOT EXISTS user_to_school
(
  user_id       NUMERIC   NOT NULL,
  school_id     NUMERIC   NOT NULL,
  from_year     INTEGER   NULL,
  to_year       INTEGER   NULL,
  insert_time   TIMESTAMP NOT NULL,
  PRIMARY KEY   (user_id,school_id,from_year,to_year)
);

CREATE TABLE IF NOT EXISTS user_to_company
(
  user_id       NUMERIC   NOT NULL,
  company_id    NUMERIC   NOT NULL,
  title         TEXT      NULL,
  from_year     INTEGER   NULL,
  to_year       INTEGER   NULL,
  insert_time   TIMESTAMP NOT NULL,
  PRIMARY KEY   (user_id,company_id,from_year,to_year)
);

CREATE TABLE IF NOT EXISTS user_to_groups
(
  user_id       NUMERIC   NOT NULL,
  group_name    TEXT      NOT NULL,
  status        BOOLEAN   NOT NULL,
  PRIMARY KEY   (user_id,group_name)
);