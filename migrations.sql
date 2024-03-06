-- https://app.sqldbm.com/PostgreSQL/DatabaseExplorer/p282867/#
CREATE TABLE IF NOT EXISTS users (
    "id" varchar(50) PRIMARY KEY NOT NULL,
    "email" varchar(100) NOT NULL,
    "phone" varchar(50) NULL,
    "phone_country" varchar(50) NULL,
    "password" varchar(255) NULL,
    "email_verify" boolean NOT NULL,
    "phone_verify" boolean NOT NULL,
    "avatar" varchar(255) NOT NULL DEFAULT "https://i1.sndcdn.com/artworks-DMKEsjVymB5A2teD-yr6bng-t240x240.jpg",
    "display_name" varchar(50) NULL,
    "username" varchar(30) NOT NULL,
    "github_connect" varchar(50) NULL,
    "google_connect" varchar(50) NULL,
    "email_2fa" boolean NOT NULL,
    "phone_2fa" boolean NOT NULL,
    "totp_2fa" boolean NOT NULL,
    "totp" char(255) NULL,
    "default_2fa" integer NOT NULL,
    "create_at" timestamp NOT NULL,
    "update_at" timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS user_profile (
    "id" varchar(50) NOT NULL,
    "bio" text NOT NULL,
    "public_email" varchar(100) NOT NULL,
    "pronouns" varchar(50) NOT NULL,
    "related_links" varchar(255) [] NOT NULL
);

CREATE INDEX IF NOT EXISTS user_profile_FK_1 ON user_profile ("id");

CREATE TABLE IF NOT EXISTS tokens (
    "id" varchar(50) PRIMARY KEY NOT NULL,
    "used_time" int NOT NULL,
    "user_agent" text NOT NULL,
    "ip" varchar(50) NOT NULL,
    "create_at" timestamp NOT NULL,
    "update_at" timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS questions (
    "id" varchar(50) PRIMARY KEY NOT NULL,
    "questioner_id" varchar(50) NOT NULL,
    "title" varchar(255) NOT NULL,
    "tags_name" varchar(255) [] NOT NULL,
    "tags_version" varchar(255) [] NOT NULL,
    "content" text NOT NULL,
    "views" integer NOT NULL,
    "create_at" timestamp NOT NULL,
    "update_at" timestamp NOT NULL,
    CONSTRAINT questions_user_id_FK_1 FOREIGN KEY ("questioner_id") REFERENCES users ("id")
);

CREATE TABLE IF NOT EXISTS tags_list (
    "tag_name" varchar(50) PRIMARY KEY NOT NULL,
    "tag_introduce" text NOT NULL,
    "create_at" timestamp NOT NULL,
    "update_at" timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS question_votes (
    "question_id" varchar(50) NOT NULL,
    "voter_id" varchar(50) NOT NULL,
    "vote" integer NOT NULL,
    PRIMARY KEY (question_id, voter_id),
    CONSTRAINT question_votes_FK_1 FOREIGN KEY ("question_id") REFERENCES questions ("id")
);

-- cus we need use question_id to get how many user vote down or up
CREATE INDEX IF NOT EXISTS question_votes_FK_1 ON question_votes ("question_id");

CREATE TABLE IF NOT EXISTS question_comments (
    "id" varchar(50) PRIMARY KEY NOT NULL,
    "question_id" varchar(50) NOT NULL,
    "commenter_id" varchar(50) NOT NULL,
    "content" text NOT NULL,
    "reply" varchar(50) NOT NULL,
    "update_at" timestamp NOT NULL,
    "create_at" timestamp NOT NULL,
    CONSTRAINT question_comments_FK_1 FOREIGN KEY ("question_id") REFERENCES questions ("id"),
    CONSTRAINT question_comments_FK_2 FOREIGN KEY ("commenter_id") REFERENCES users ("id")
);

-- cus we need use question_id to get all comment to show on the fornt side
CREATE INDEX IF NOT EXISTS question_comments_FK_1 ON question_comments ("question_id");

CREATE TABLE IF NOT EXISTS question_comment_votes (
    "comment_id" varchar(50) NOT NULL,
    "voter_id" varchar(50) NOT NULL,
    "vote" integer NOT NULL,
    CONSTRAINT question_comment_votes_PK_1 PRIMARY KEY ("comment_id", "voter_id"),
    CONSTRAINT question_comment_votes_FK_1 FOREIGN KEY ("comment_id") REFERENCES question_comments ("id")
);

-- cus need use comment_id to get how many user vote up or donw so need create a index
CREATE INDEX IF NOT EXISTS question_comment_votes_FK_1 ON question_comment_votes ("comment_id");

CREATE TABLE IF NOT EXISTS question_answers (
    "id" varchar(50) PRIMARY KEY NOT NULL,
    "question_id" varchar(50) NOT NULL,
    "user_id" varchar(50) NOT NULL,
    "content" text NOT NULL,
    "create_at" timestamp NOT NULL,
    "update_at" timestamp NOT NULL,
    CONSTRAINT question_answers_FK_1 FOREIGN KEY ("question_id") REFERENCES questions ("id"),
    CONSTRAINT question_answers_FK_2 FOREIGN KEY ("user_id") REFERENCES users ("id")
);

-- cus need user question_id to get all question_anwsers
CREATE INDEX IF NOT EXISTS question_answers_FK_1 ON question_answers ("question_id");

CREATE TABLE IF NOT EXISTS question_answer_votes (
    "answer_id" varchar(50) NOT NULL,
    "voter_id" varchar(50) NOT NULL,
    "vote" integer NOT NULL,
    CONSTRAINT question_answer_votes_PK_1 PRIMARY KEY ("answer_id", "voter_id"),
    CONSTRAINT question_answer_votes_FK_1 FOREIGN KEY ("answer_id") REFERENCES question_answers ("id")
);

-- cus we need use answer_id to get how many user vote up or down so need create a index
CREATE INDEX IF NOT EXISTS question_answer_votes_FK_1 ON question_answer_votes ("answer_id");

CREATE TABLE IF NOT EXISTS question_chat (
    "id" varchar(50) PRIMARY KEY NOT NULL,
    "question_id" varchar(50) NOT NULL,
    "reply" varchar(50),
    "author" varchar(50) NOT NULL,
    "content" varchar(50) NOT NULL,
    "image" varchar(255) [] NOT NULL,
    "create_at" timestamp NOT NULL,
    "update_at" timestamp NOT NULL,
    CONSTRAINT question_chat_FK_1 FOREIGN KEY ("question_id") REFERENCES questions ("id")
);