CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    album VARCHAR(255),
    release_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);