-- Restore UNIQUE(tournament_id, place) if missing
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conrelid = 'results'::regclass
          AND conname = 'results_tournament_id_place_key'
    ) THEN
        ALTER TABLE results
            ADD CONSTRAINT results_tournament_id_place_key UNIQUE (tournament_id, place);
    END IF;
END $$;
