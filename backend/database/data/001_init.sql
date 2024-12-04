CREATE TABLE "demos" (
  "id" integer  PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "user_id" integer NOT NULL,
  "link" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "upvotes" integer NOT NULL,
  "downvotes" integer NOT NULL,
  "topic_id" integer NOT NULL
);

CREATE TABLE "demo_access" (
  "demo_id" integer NOT NULL,
  "user_id" integer NOT NULL
);

CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "username" varchar(255) NOT NULL,
  "display_name" varchar(255),
  "email" varchar(255) NOT NULL UNIQUE,
  "password" varchar(255) NOT NULL,
  "role_id" integer NOT NULL,
  "created_at" timestamp NOT NULL,
  "karma" integer NOT NULL
);

CREATE TABLE "roles" (
  "id" integer PRIMARY KEY,
  "name" varchar(255) NOT NULL
);

CREATE TABLE "forums" (
  "id" integer PRIMARY KEY,
  "name" varchar(255) NOT NULL
);

CREATE TABLE "topics" (
  "id" integer PRIMARY KEY,
  "title" varchar(255), NOT NULL
  "user_id" integer NOT NULL,
  "forum_id" integer NOT NULL,
  "tag" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL,
  "last_update" timestamp NOT NULL,
  "total_upvotes" integer NOT NULL
  "total_downvotes" integer NOT NULL
);

CREATE TABLE "messages" (
  "id" integer PRIMARY KEY,
  "topic_id" integer NOT NULL,
  "user_id" integer NOT NULL,
  "title" varchar(255) NOT NULL,
  "body" text NOT NULL,
  "upvotes" integer NOT NULL
  "downvotes" integer NOT NULL
);

CREATE TABLE "assets" (
  "id" integer PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "link" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL
);

ALTER TABLE "users" ADD FOREIGN KEY ("id") REFERENCES "demos" ("user_id");

ALTER TABLE "topics" ADD FOREIGN KEY ("id") REFERENCES "demos" ("topic_id");

ALTER TABLE "users" ADD FOREIGN KEY ("id") REFERENCES "demo_access" ("user_id");

ALTER TABLE "demos" ADD FOREIGN KEY ("id") REFERENCES "demo_access" ("demo_id");

ALTER TABLE "roles" ADD FOREIGN KEY ("id") REFERENCES "users" ("role_id");

ALTER TABLE "topics" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "topics" ADD FOREIGN KEY ("forum_id") REFERENCES "forums" ("id");

ALTER TABLE "messages" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "messages" ADD FOREIGN KEY ("topic_id") REFERENCES "topics" ("id");

---- create above / drop below ----

DROP TABLE "demos";
DROP TABLE "demo_access";
DROP TABLE "users";
DROP TABLE "roles";
DROP TABLE "forums";
DROP TABLE "topics";
DROP TABLE "messages";
DROP TABLE "assets";
