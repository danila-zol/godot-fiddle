CREATE SCHEMA asset;

CREATE TABLE asset.assets (
  "id" varchar(64) PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" varchar NOT NULL,
  "link" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL
);

---- create above / drop below ----

DROP SCHEMA asset CASCADE;
