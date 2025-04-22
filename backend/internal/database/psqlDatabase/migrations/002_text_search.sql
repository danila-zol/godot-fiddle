ALTER TABLE demo.demos ADD COLUMN ts tsvector GENERATED ALWAYS AS (
	setweight(to_tsvector('english', "title"), 'A') ||
	setweight(to_tsvector('english', COALESCE("description", '')), 'B')
) STORED;
CREATE INDEX demo_gin_index_ts ON demo.demos USING GIN (ts);

---- create above / drop below ----

ALTER TABLE demo.demos DROP COLUMN ts;
DROP INDEX demo_gin_index_ts;
