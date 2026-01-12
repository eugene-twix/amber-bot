-- Add recorded_by column and rename created_at to recorded_at in results
ALTER TABLE results ADD COLUMN IF NOT EXISTS recorded_by BIGINT;

-- Rename created_at to recorded_at if exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'results' AND column_name = 'created_at') THEN
        ALTER TABLE results RENAME COLUMN created_at TO recorded_at;
    ELSIF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'results' AND column_name = 'recorded_at') THEN
        ALTER TABLE results ADD COLUMN recorded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;
    END IF;
END $$;

-- Rename created_at to joined_at in members (or add if missing)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'members' AND column_name = 'created_at') THEN
        ALTER TABLE members RENAME COLUMN created_at TO joined_at;
    ELSIF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'members' AND column_name = 'joined_at') THEN
        ALTER TABLE members ADD COLUMN joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;
    END IF;
END $$;
