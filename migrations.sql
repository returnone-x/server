CREATE TABLE users (
    "id" varchar(50) PRIMARY KEY NOT NULL,
    email varchar(100) NOT NULL,
    phone varchar(50) NULL,
    phone_country varchar(50) NULL,
    password varchar(255) NULL,
    email_verify boolean NOT NULL,
    phone_verify boolean NOT NULL,
    avatar varchar(100) NOT NULL DEFAULT "https://i1.sndcdn.com/artworks-DMKEsjVymB5A2teD-yr6bng-t240x240.jpg",
    display_name varchar(50) NULL,
    user_name varchar(30) NOT NULL,
    github_connect varchar(50) NULL,
    google_connect varchar(50) NULL,
    email_2fa boolean NOT NULL,
    phone_2fa boolean NOT NULL,
    totp_2fa boolean NOT NULL,
    totp char(255) NULL,
    default_2fa integer NOT NULL,
    create_at timestamp NOT NULL,
    update_at timestamp NOT NULL
);

CREATE TABLE tokens (
    "id" varchar(50) PRIMARY KEY NOT NULL,
    used_time int NOT NULL,
    user_agent text NOT NULL,
    ip varchar(50) NOT NULL,
    create_at timestamp NOT NULL,
    update_at timestamp NOT NULL
);

CREATE TABLE questions (
    "id" varchar(50) PRIMARY KEY NOT NULL,
    questioner_id varchar(50) NOT NULL,
    title varchar(255) NOT NULL,
    tags jsonb [] NOT NULL,
    content text NOT NULL,
    views integer NOT NULL,
    create_at timestamp NOT NULL,
    update_at timestamp NOT NULL,
    CONSTRAINT questions_user_id_FK_1 FOREIGN KEY (questioner_id) REFERENCES users ("id")
);

CREATE INDEX questions_user_id_FK_1 ON questions (questioner_id);