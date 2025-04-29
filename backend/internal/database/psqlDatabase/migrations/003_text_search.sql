CREATE OR REPLACE FUNCTION to_tsvector_multilang(varchar) RETURNS tsvector AS $$
BEGIN
	RETURN 
	to_tsvector('english', $1) || 
	to_tsvector('russian', $1);
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE OR REPLACE FUNCTION to_tsquery_multilang(varchar) RETURNS tsquery AS $$
BEGIN
	RETURN
	websearch_to_tsquery('english', $1) || 
	websearch_to_tsquery('russian', $1);
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Demo
ALTER TABLE demo.demos ADD COLUMN demo_ts tsvector GENERATED ALWAYS AS (
	setweight(to_tsvector_multilang("title"), 'A') ||
	setweight(to_tsvector_multilang(COALESCE("description", '')), 'B')
) STORED;
CREATE INDEX demo_gin_index_ts ON demo.demos USING GIN (demo_ts);

-- Asset
ALTER TABLE asset.assets ADD COLUMN asset_ts tsvector GENERATED ALWAYS AS (
	setweight(to_tsvector_multilang("name"), 'A') ||
	setweight(to_tsvector_multilang(COALESCE("description", '')), 'B')
) STORED;
CREATE INDEX asset_gin_index_ts ON asset.assets USING GIN (asset_ts);

-- Thread
ALTER TABLE forum.threads ADD COLUMN thread_ts tsvector GENERATED ALWAYS AS (
	setweight(to_tsvector_multilang("title"), 'A')
) STORED;
CREATE INDEX thread_gin_index_ts ON forum.threads USING GIN (thread_ts);

-- Message
ALTER TABLE forum.messages ADD COLUMN message_ts tsvector GENERATED ALWAYS AS (
	setweight(to_tsvector_multilang("title"), 'A') ||
	setweight(to_tsvector_multilang(COALESCE("body", '')), 'B')
) STORED;
CREATE INDEX message_gin_index_ts ON forum.messages USING GIN (message_ts);

CREATE COLLATION case_insensitive (provider = icu, locale = 'und-u-ks-level2', deterministic = false);

---- create above / drop below ----

ALTER TABLE demo.demos DROP COLUMN IF EXISTS demo_ts;
ALTER TABLE asset.assets DROP COLUMN IF EXISTS asset_ts;
ALTER TABLE forum.threads DROP COLUMN IF EXISTS thread_ts;
ALTER TABLE forum.messages DROP COLUMN IF EXISTS message_ts;
DROP INDEX IF EXISTS demo_gin_index_ts;
DROP INDEX IF EXISTS asset_gin_index_ts;
DROP INDEX IF EXISTS thread_gin_index_ts;
DROP INDEX IF EXISTS message_gin_index_ts;
DROP FUNCTION IF EXISTS to_tsvector_multilang;
DROP FUNCTION IF EXISTS to_tsquery_multilang;
DROP COLLATION IF EXISTS case_insensitive;
