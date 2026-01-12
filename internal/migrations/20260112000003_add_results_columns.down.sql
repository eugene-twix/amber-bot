-- Remove recorded_by and recorded_at columns from results
ALTER TABLE results DROP COLUMN IF EXISTS recorded_by;
ALTER TABLE results DROP COLUMN IF EXISTS recorded_at;

-- Rename joined_at back to created_at in members
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'members' AND column_name = 'joined_at') THEN
        ALTER TABLE members RENAME COLUMN joined_at TO created_at;
    END IF;
END $$;
