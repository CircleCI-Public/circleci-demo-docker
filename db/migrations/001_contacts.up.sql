CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    email varchar(255) UNIQUE NOT NULL,
    name varchar(255) NOT NULL
);
