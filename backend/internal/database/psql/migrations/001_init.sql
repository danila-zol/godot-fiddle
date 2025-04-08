CREATE SCHEMA "user";

CREATE TABLE "user".roles (
	"id" varchar(64) PRIMARY KEY,
	"name" varchar(255) NOT NULL,
	"permissions" []varchar(64)
);

CREATE TABLE "user".users (
	"id" varchar(64) PRIMARY KEY,
	"username" varchar(255) NOT NULL UNIQUE,
	"displayName" varchar(255),
	"email" varchar(255) NOT NULL UNIQUE,
	"password" varchar(255) NOT NULL,
	"roleID" varchar(64) NOT NULL REFERENCES "user".roles (id) ON DELETE RESTRICT,
	"createdAt" timestamp NOT NULL,
	"karma" integer NOT NULL
);

CREATE TABLE "user".sessions (
	"id" varchar(64) PRIMARY KEY,
	"userID" varchar(64) NOT NULL REFERENCES "user".users (id) ON DELETE RESTRICT
);

CREATE SCHEMA demo;

CREATE TABLE demo.demos (
	"id" varchar(64)  PRIMARY KEY,
	"name" varchar(255) NOT NULL,
	"description" varchar NOT NULL,
	"userID" varchar(64) NOT NULL,
	"link" varchar(255) NOT NULL,
	"createdAt" timestamp NOT NULL,
	"updatedAt" timestamp NOT NULL,
	"upvotes" integer NOT NULL,
	"downvotes" integer NOT NULL,
	"threadID" varchar(64) NOT NULL     -- Links to a thread in the forums
);

CREATE SCHEMA forum;

CREATE TABLE forum.topics (
	"id" varchar(64) PRIMARY KEY,
	"name" varchar(255) NOT NULL
);

CREATE TABLE forum.threads (
	"id" varchar(64) PRIMARY KEY,
	"title" varchar(255) NOT NULL,
	"userID" varchar(64) NOT NULL,
	"topicID" varchar(64) NOT NULL REFERENCES forum.topics (id) ON DELETE CASCADE,
	"tags" []varchar(255) NOT NULL,
	"createdAt" timestamp NOT NULL,
	"lastUpdate" timestamp NOT NULL,
	"totalUpvotes" integer NOT NULL,
	"totalDownvotes" integer NOT NULL
);

CREATE TABLE forum.messages (
	"id" varchar(64) PRIMARY KEY,
	"threadID" varchar(64) NOT NULL REFERENCES forum.threads (id) ON DELETE CASCADE,
	"userID" varchar(64) NOT NULL,
	"title" varchar(255) NOT NULL,
	"body" varchar NOT NULL,
	"tags" []varchar(255) NOT NULL,
	"createdAt" timestamp NOT NULL,
	"updatedAt" timestamp NOT NULL,
	"upvotes" integer NOT NULL,
	"downvotes" integer NOT NULL
);

CREATE SCHEMA asset;

CREATE TABLE asset.assets (
  "id" varchar(64) PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" varchar NOT NULL,
  "link" varchar(255) NOT NULL,
  "createdAt" timestamp NOT NULL
);

---- create above / drop below ----

DROP SCHEMA "user" CASCADE;
DROP SCHEMA "demo" CASCADE;
DROP SCHEMA "forum" CASCADE;
DROP SCHEMA "asset" CASCADE;
