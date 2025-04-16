CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SCHEMA IF NOT EXISTS "user";

CREATE TABLE "user".roles (
	"id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	"name" varchar(255) NOT NULL
	-- "permissions" varchar(64)[]
);

CREATE TABLE "user".users (
	"id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	"username" varchar(255) NOT NULL UNIQUE,
	"display_name" varchar(255),
	"email" varchar(255) NOT NULL UNIQUE,
	"password" varchar(255) NOT NULL,
	"verified" boolean DEFAULT FALSE NOT NULL,
	"role_id" uuid NOT NULL REFERENCES "user".roles (id) ON DELETE RESTRICT,
	"created_at" timestamp DEFAULT NOW() NOT NULL,
	"karma" integer DEFAULT 0 NOT NULL
);

CREATE TABLE "user".sessions (
	"id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	"user_id" uuid NOT NULL REFERENCES "user".users (id) ON DELETE CASCADE
);

CREATE SCHEMA IF NOT EXISTS forum;

CREATE TABLE forum.topics (
	"id"  integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	"name" varchar(255) NOT NULL
);

CREATE TABLE forum.threads (
	"id" integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	"title" varchar(255) NOT NULL,
	"user_id" uuid NOT NULL,
	"topic_id" integer NOT NULL REFERENCES forum.topics (id) ON DELETE CASCADE,
	"tags" varchar(255)[],
	"created_at" timestamp DEFAULT NOW() NOT NULL,
	"updated_at" timestamp DEFAULT NOW() NOT NULL,
	"upvotes" integer NOT NULL,
	"downvotes" integer NOT NULL
);

CREATE TABLE forum.messages (
	"id" integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	"thread_id" integer NOT NULL REFERENCES forum.threads (id) ON DELETE CASCADE,
	"user_id" uuid NOT NULL,
	"title" varchar(255) NOT NULL,
	"body" varchar NOT NULL,
	"tags" varchar(255)[],
	"created_at" timestamp DEFAULT NOW() NOT NULL,
	"updated_at" timestamp DEFAULT NOW() NOT NULL,
	"upvotes" integer NOT NULL,
	"downvotes" integer NOT NULL
);

CREATE SCHEMA IF NOT EXISTS demo;

CREATE TABLE demo.demos (
	"id"  integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	"title" varchar(255) NOT NULL,
	"description" varchar,
	"tags" varchar(255)[],
	"link" varchar(255) NOT NULL,
	"user_id" uuid NOT NULL,
	"created_at" timestamp DEFAULT NOW() NOT NULL,
	"updated_at" timestamp DEFAULT NOW() NOT NULL,
	"upvotes" integer NOT NULL,
	"downvotes" integer NOT NULL,
	"thread_id" integer NOT NULL REFERENCES forum.threads (id) ON DELETE CASCADE
);

CREATE SCHEMA IF NOT EXISTS asset;

CREATE TABLE asset.assets (
  "id" integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  "name" varchar(255) NOT NULL,
  "description" varchar,
  "link" varchar(255) NOT NULL,
  "created_at" timestamp DEFAULT NOW() NOT NULL
);

---- create above / drop below ----

DROP SCHEMA IF EXISTS "user" CASCADE;
DROP SCHEMA IF EXISTS "demo" CASCADE;
DROP SCHEMA IF EXISTS "forum" CASCADE;
DROP SCHEMA IF EXISTS "asset" CASCADE;
DROP EXTENSION IF EXISTS "uuid-ossp";
