CREATE TABLE IF NOT EXISTS linkedin_user
(
    user_id     TEXT      NOT NULL,
    first_name  TEXT      NOT NULL,
    last_name   TEXT      NOT NULL,
    url         TEXT      NULL,
    picture     TEXT      NULL,
    insert_time TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS linkedin_user_token
(
    user_id TEXT NOT NULL,
    token   TEXT NOT NULL,
    insert_time TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id)
);


CREATE TABLE IF NOT EXISTS school
(
    school_id      SERIAL,
    school_name    TEXT      NOT NULL,
    degree         TEXT      NULL,
    field_of_study TEXT      NULL,
    insert_time    TIMESTAMP NOT NULL,
    PRIMARY KEY (school_id)
);

CREATE TABLE IF NOT EXISTS company
(
    company_id   SERIAL,
    company_name TEXT      NOT NULL,
    location     TEXT      NULL,
    insert_time  TIMESTAMP NOT NULL,
    PRIMARY KEY (company_id)
);

CREATE TABLE IF NOT EXISTS user_to_school
(
    user_id     TEXT   NOT NULL,
    school_id   NUMERIC   NOT NULL,
    from_year   INTEGER   NULL,
    to_year     INTEGER   NULL,
    status      BOOLEAN NOT NULL,
    insert_time TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, school_id, from_year, to_year)
);

CREATE TABLE IF NOT EXISTS user_to_company
(
    user_id     TEXT   NOT NULL,
    company_id  NUMERIC   NOT NULL,
    title       TEXT      NULL,
    from_year   INTEGER   NULL,
    to_year     INTEGER   NULL,
    status      BOOLEAN NOT NULL,
    insert_time TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, company_id, from_year, to_year)
);

CREATE TABLE IF NOT EXISTS user_to_groups
(
    user_id      TEXT NOT NULL,
    group_name   TEXT    NOT NULL,
    status       BOOLEAN NOT NULL,
    group_source TEXT NOT NULL,
    PRIMARY KEY  (user_id, group_name)
);