ALTER TABLE asset.assets ADD COLUMN "version" INTEGER NOT NULL DEFAULT 1;

CREATE FUNCTION increment_version()
	RETURNS TRIGGER AS
$func$
BEGIN
	NEW.version := OLD.version + 1;
	RETURN NEW;
END;
$func$ LANGUAGE plpgsql;

CREATE TRIGGER increment_version_on_update
	BEFORE UPDATE ON asset.assets
	FOR EACH ROW
	EXECUTE FUNCTION increment_version();

---- create above / drop below ----

ALTER TABLE asset.assets DROP COLUMN IF EXISTS "version";
DROP TRIGGER IF EXISTS increment_version_on_update ON asset.assets; 
DROP FUNCTION IF EXISTS increment_version;
