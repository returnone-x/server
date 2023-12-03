CREATE TABLE users
(
 "id"           varchar(50) NOT NULL,
 email          varchar(100) NOT NULL,
 phone          varchar(50) NULL,
 phone_country  varchar(50) NULL,
 password       varchar(255) NULL,
 email_verify   boolean NOT NULL,
 phone_verify   boolean NOT NULL,
 avatar         varchar(100) NULL,
 display_name   varchar(50) NULL,
 user_name      varchar(30) NOT NULL,
 github_connect varchar(100) NULL,
 google_connect varchar(100) NULL,
 email_2fa      boolean NOT NULL,
 phone_2fa      boolean NOT NULL,
 totp_2fa       boolean NOT NULL,
 totp           char(255) NULL,
 default_2fa    integer NOT NULL,
 create_at      timestamp NOT NULL,
 update_at      timestamp NOT NULL,
 CONSTRAINT PK_1 PRIMARY KEY ( "id" )
);