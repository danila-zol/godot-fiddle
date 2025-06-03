ALTER TABLE asset.assets ADD COLUMN "version" INTEGER NOT NULL DEFAULT 1;
ALTER TABLE forum.topics ADD COLUMN "version" INTEGER NOT NULL DEFAULT 1;

CREATE FUNCTION increment_version()
	RETURNS TRIGGER AS
$func$
BEGIN
	NEW.version := OLD.version + 1;
	RETURN NEW;
END;
$func$ LANGUAGE plpgsql;

CREATE TRIGGER increment_asset_version_on_update
	BEFORE UPDATE ON asset.assets
	FOR EACH ROW
		WHEN ((OLD.name IS DISTINCT FROM NEW.name) 
			OR (OLD.description IS DISTINCT FROM NEW.description)
			OR (OLD.tags IS DISTINCT FROM NEW.tags))
	EXECUTE FUNCTION increment_version();

CREATE TRIGGER increment_topic_version_on_update
	BEFORE UPDATE ON forum.topics
	FOR EACH ROW
	EXECUTE FUNCTION increment_version();

---- create above / drop below ----

ALTER TABLE asset.assets DROP COLUMN IF EXISTS "version";
ALTER TABLE forum.topics DROP COLUMN IF EXISTS "version";
DROP TRIGGER IF EXISTS increment_asset_version_on_update ON asset.assets; 
DROP TRIGGER IF EXISTS increment_topic_version_on_update ON forum.topics; 
DROP FUNCTION IF EXISTS increment_version;
