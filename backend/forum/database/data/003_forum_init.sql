CREATE SCHEMA forum;

CREATE TABLE forum.topics (
  "id" varchar(64) PRIMARY KEY,
  "name" varchar(255) NOT NULL
);

CREATE TABLE forum.threads (
  "id" varchar(64) PRIMARY KEY,
  "title" varchar(255) NOT NULL,
  "user_id" varchar(64) NOT NULL,
  "topic_id" varchar(64) NOT NULL REFERENCES forum.topics (id) ON DELETE CASCADE,
  "tag" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL,
  "last_update" timestamp NOT NULL,
  "total_upvotes" integer NOT NULL,
  "total_downvotes" integer NOT NULL
);

CREATE TABLE forum.messages (
  "id" varchar(64) PRIMARY KEY,
  "thread_id" varchar(64) NOT NULL REFERENCES forum.threads (id) ON DELETE CASCADE,
  "user_id" varchar(64) NOT NULL,
  "title" varchar(255) NOT NULL,
  "body" varchar NOT NULL,
  "tag" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "upvotes" integer NOT NULL,
  "downvotes" integer NOT NULL
);

---- create above / drop below ----

DROP SCHEMA forum CASCADE;
