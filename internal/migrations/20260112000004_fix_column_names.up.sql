-- Fix results: rename created_at to recorded_at if needed
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'results' AND column_name = 'created_at') THEN
        ALTER TABLE results RENAME COLUMN created_at TO recorded_at;
    END IF;
END $$;
