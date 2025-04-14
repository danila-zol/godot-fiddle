CREATE SCHEMA "user";

CREATE TABLE "user".roles (
	"id" varchar(64) PRIMARY KEY,
	"name" varchar(255) NOT NULL
	-- "permissions" varchar(64)[]
);

CREATE TABLE "user".users (
	"id" varchar(64) PRIMARY KEY,
	"username" varchar(255) NOT NULL UNIQUE,
	"displayName" varchar(255),
	"email" varchar(255) NOT NULL UNIQUE,
	"password" varchar(255) NOT NULL,
	"verified" boolean NOT NULL DEFAULT false,
	"roleID" varchar(64) NOT NULL REFERENCES "user".roles (id) ON DELETE RESTRICT,
	"createdAt" timestamp NOT NULL,
	"karma" integer NOT NULL
);

CREATE TABLE "user".sessions (
	"id" varchar(64) PRIMARY KEY,
	"userID" varchar(64) NOT NULL REFERENCES "user".users (id) ON DELETE RESTRICT
);

CREATE SCHEMA forum;

CREATE TABLE forum.topics (
	"id"  integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	"name" varchar(255) NOT NULL
);

CREATE TABLE forum.threads (
	"id" integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	"title" varchar(255) NOT NULL,
	"userID" varchar(64) NOT NULL,
	"topicID" integer NOT NULL REFERENCES forum.topics (id) ON DELETE CASCADE,
	"tags" varchar(255)[],
	"createdAt" timestamp NOT NULL,
	"lastUpdate" timestamp NOT NULL,
	"totalUpvotes" integer NOT NULL,
	"totalDownvotes" integer NOT NULL
);

CREATE TABLE forum.messages (
	"id" integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	"threadID" integer NOT NULL REFERENCES forum.threads (id) ON DELETE CASCADE,
	"userID" varchar(64) NOT NULL,
	"title" varchar(255) NOT NULL,
	"body" varchar NOT NULL,
	"tags" varchar(255)[],
	"createdAt" timestamp NOT NULL,
	"updatedAt" timestamp NOT NULL,
	"upvotes" integer NOT NULL,
	"downvotes" integer NOT NULL
);

CREATE SCHEMA demo;

CREATE TABLE demo.demos (
	"id"  integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	"title" varchar(255) NOT NULL,
	"description" varchar NOT NULL,
	"tags" varchar(255)[],
	"link" varchar(255) NOT NULL,
	"userID" varchar(64) NOT NULL,
	"createdAt" timestamp NOT NULL,
	"updatedAt" timestamp NOT NULL,
	"upvotes" integer NOT NULL,
	"downvotes" integer NOT NULL,
	"threadID" integer NOT NULL REFERENCES forum.threads (id) ON DELETE CASCADE
);

CREATE SCHEMA asset;

CREATE TABLE asset.assets (
  "id" integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  "name" varchar(255) NOT NULL,
  "description" varchar,
  "link" varchar(255) NOT NULL,
  "createdAt" timestamp NOT NULL
);

---- create above / drop below ----

DROP SCHEMA "user" CASCADE;
DROP SCHEMA "demo" CASCADE;
DROP SCHEMA "forum" CASCADE;
DROP SCHEMA "asset" CASCADE;
