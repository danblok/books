CREATE TABLE IF NOT EXISTS Books (
    id uuid PRIMARY KEY,
    name text UNIQUE,
    author text,
    price integer
);
