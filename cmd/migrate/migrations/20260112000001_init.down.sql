-- Drop tables in reverse order (respecting foreign keys)
DROP TABLE IF EXISTS results;
DROP TABLE IF EXISTS members;
DROP TABLE IF EXISTS tournaments;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;

-- Drop enum type
DROP TYPE IF EXISTS user_role;
