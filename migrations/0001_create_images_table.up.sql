CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    path TEXT UNIQUE,
    dimensions TEXT,
    camera_model TEXT,
    location TEXT,
    format TEXT,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
