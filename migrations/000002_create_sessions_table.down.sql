-- Remove index on sessions table --
DROP INDEX sessions_expiry_idx ON sessions;
-- Remove sessions table --
DROP TABLE sessions;
