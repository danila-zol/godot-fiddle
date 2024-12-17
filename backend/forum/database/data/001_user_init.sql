CREATE SCHEMA "user";

CREATE TABLE "user".roles (
  "id" varchar(64) PRIMARY KEY,
  "name" varchar(255) NOT NULL
);

CREATE TABLE "user".users (
  "id" varchar(64) PRIMARY KEY,
  "username" varchar(255) NOT NULL UNIQUE,
  "display_name" varchar(255),
  "email" varchar(255) NOT NULL UNIQUE,
  "password" varchar(255) NOT NULL,
  "role_id" varchar(64) NOT NULL REFERENCES "user".roles (id) ON DELETE RESTRICT,
  "created_at" timestamp NOT NULL,
  "karma" integer NOT NULL
);

---- create above / drop below ----

DROP SCHEMA "user" CASCADE;
