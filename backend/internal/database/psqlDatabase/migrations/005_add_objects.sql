ALTER TABLE demo.demos ADD COLUMN object_key VARCHAR(255) GENERATED ALWAYS AS ('demo-' || demo.demos.id) STORED;
ALTER TABLE demo.demos ADD COLUMN thumbnail_key VARCHAR(255) GENERATED ALWAYS AS ('demo-thumnail-' || demo.demos.id) STORED;

ALTER TABLE asset.assets ADD COLUMN object_key VARCHAR(255) GENERATED ALWAYS AS ('asset-' || asset.assets.id) STORED;
ALTER TABLE asset.assets ADD COLUMN thumbnail_key VARCHAR(255) GENERATED ALWAYS AS ('asset-thumbnail-' || asset.assets.id) STORED;

-- User avatars keys are their unique usernames

ALTER TABLE demo.demos DROP COLUMN IF EXISTS "link";
ALTER TABLE asset.assets DROP COLUMN IF EXISTS "link";

---- create above / drop below ----

ALTER TABLE demo.demos DROP COLUMN IF EXISTS object_key;
ALTER TABLE demo.demos DROP COLUMN IF EXISTS thumbnail_key;

ALTER TABLE asset.assets DROP COLUMN IF EXISTS object_key;
ALTER TABLE asset.assets DROP COLUMN IF EXISTS thumbnail_key;

ALTER TABLE demo.demos ADD COLUMN IF NOT EXISTS "link" VARCHAR(255) NOT NULL;
ALTER TABLE asset.assets ADD COLUMN IF NOT EXISTS "link" VARCHAR(255) NOT NULL;
