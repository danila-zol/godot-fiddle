CREATE TABLE "roles" (
  "id" integer PRIMARY KEY,
  "name" varchar(255) NOT NULL
);

CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "username" varchar(255) NOT NULL,
  "display_name" varchar(255),
  "email" varchar(255) NOT NULL UNIQUE,
  "password" varchar(255) NOT NULL,
  "role_id" integer NOT NULL REFERENCES "roles" ("id"),
  "created_at" timestamp NOT NULL,
  "karma" integer NOT NULL
);

CREATE TABLE "forums" (
  "id" integer PRIMARY KEY,
  "name" varchar(255) NOT NULL
);

CREATE TABLE "topics" (
  "id" integer PRIMARY KEY,
  "title" varchar(255) NOT NULL,
  "user_id" integer NOT NULL REFERENCES "users" ("id"),
  "forum_id" integer NOT NULL REFERENCES "forums" ("id"),
  "tag" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL,
  "last_update" timestamp NOT NULL,
  "total_upvotes" integer NOT NULL,
  "total_downvotes" integer NOT NULL
);

CREATE TABLE "demos" (
  "id" integer  PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" varchar NOT NULL,
  "user_id" integer NOT NULL REFERENCES "users" ("id"),
  "link" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "upvotes" integer NOT NULL,
  "downvotes" integer NOT NULL,
  "topic_id" integer NOT NULL REFERENCES "topics" ("id")
);

CREATE TABLE "demo_access" (
  "demo_id" integer NOT NULL REFERENCES "demos" ("id"),
  "user_id" integer NOT NULL REFERENCES "users" ("id")
);

CREATE TABLE "messages" (
  "id" integer PRIMARY KEY,
  "topic_id" integer NOT NULL REFERENCES "topics" ("id"),
  "user_id" integer NOT NULL REFERENCES "users" ("id"),
  "title" varchar(255) NOT NULL,
  "body" varchar NOT NULL,
  "upvotes" integer NOT NULL,
  "downvotes" integer NOT NULL
);

CREATE TABLE "assets" (
  "id" integer PRIMARY KEY,
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
