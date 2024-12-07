CREATE TABLE "roles" (
  "id" varchar(64) PRIMARY KEY,
  "name" varchar(255) NOT NULL
);

CREATE TABLE "users" (
  "id" varchar(64) PRIMARY KEY,
  "username" varchar(255) NOT NULL UNIQUE,
  "display_name" varchar(255),
  "email" varchar(255) NOT NULL UNIQUE,
  "password" varchar(255) NOT NULL,
  "role_id" varchar(64) NOT NULL REFERENCES "roles" ("id") ON DELETE RESTRICT,
  "created_at" timestamp NOT NULL,
  "karma" integer NOT NULL
);

CREATE TABLE "forums" (
  "id" varchar(64) PRIMARY KEY,
  "name" varchar(255) NOT NULL
);

CREATE TABLE "topics" (
  "id" varchar(64) PRIMARY KEY,
  "title" varchar(255) NOT NULL,
  "user_id" varchar(64) NOT NULL REFERENCES "users" ("id") ON DELETE SET NULL,
  "forum_id" varchar(64) NOT NULL REFERENCES "forums" ("id") ON DELETE CASCADE,
  "tag" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL,
  "last_update" timestamp NOT NULL,
  "total_upvotes" integer NOT NULL,
  "total_downvotes" integer NOT NULL
);

CREATE TABLE "messages" (
  "id" varchar(64) PRIMARY KEY,
  "topic_id" varchar(64) NOT NULL REFERENCES "topics" ("id") ON DELETE CASCADE,
  "user_id" varchar(64) NOT NULL REFERENCES "users" ("id") ON DELETE SET NULL,
  "title" varchar(255) NOT NULL,
  "body" varchar NOT NULL,
  "upvotes" integer NOT NULL,
  "downvotes" integer NOT NULL
);

CREATE TABLE "demos" (
  "id" varchar(64)  PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" varchar NOT NULL,
  "user_id" varchar(64) NOT NULL REFERENCES "users" ("id") ON DELETE SET NULL,
  "link" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "upvotes" integer NOT NULL,
  "downvotes" integer NOT NULL,
  "topic_id" varchar(64) NOT NULL REFERENCES "topics" ("id") ON DELETE CASCADE
);

CREATE TABLE "demo_access" (
  "demo_id" varchar(64) NOT NULL REFERENCES "demos" ("id") ON DELETE CASCADE,
  "user_id" varchar(64) NOT NULL REFERENCES "users" ("id") ON DELETE RESTRICT,
  PRIMARY KEY ("demo_id", "user_id")
);

CREATE TABLE "assets" (
  "id" varchar(64) PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" varchar NOT NULL,
  "link" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL
);

---- create above / drop below ----

DROP TABLE "demos";
DROP TABLE "demo_access";
DROP TABLE "users";
DROP TABLE "roles";
DROP TABLE "forums";
DROP TABLE "topics";
DROP TABLE "messages";
DROP TABLE "assets";
