-- Revert: rename recorded_at back to created_at
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'results' AND column_name = 'recorded_at') THEN
        ALTER TABLE results RENAME COLUMN recorded_at TO created_at;
    END IF;
END $$;
