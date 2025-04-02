CREATE TABLE plants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    species VARCHAR(255) NOT NULL,
    health INT NOT NULL CHECK (health BETWEEN 0 AND 100),
    water INT NOT NULL CHECK (water BETWEEN 0 AND 100),
    fertilizer INT NOT NULL CHECK (fertilizer BETWEEN 0 AND 100),
    planted_at TIMESTAMP NOT NULL,
    last_watered TIMESTAMP,
    last_fertilized TIMESTAMP,
    stage VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL
);