CREATE TABLE upload_links (
    id SERIAL PRIMARY KEY,
    token TEXT UNIQUE,
    expiration TIMESTAMP,
    used BOOLEAN DEFAULT FALSE
);
