CREATE SCHEMA demo;

CREATE TABLE demo.demos (
  "id" varchar(64)  PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" varchar NOT NULL,
  "user_id" varchar(64) NOT NULL,
  "link" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "upvotes" integer NOT NULL,
  "downvotes" integer NOT NULL,
  "thread_id" varchar(64) NOT NULL     -- Links to a thread in the forums
);

---- create above / drop below ----

DROP SCHEMA demo CASCADE;
