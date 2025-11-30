-- Create table for user sessions --
CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

-- Create index for sessions table --
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
