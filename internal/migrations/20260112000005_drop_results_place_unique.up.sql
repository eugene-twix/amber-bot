-- Drop UNIQUE(tournament_id, place) to allow updating places per team
ALTER TABLE results DROP CONSTRAINT IF EXISTS results_tournament_id_place_key;
