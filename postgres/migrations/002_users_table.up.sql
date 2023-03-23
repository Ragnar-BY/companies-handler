
CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
    username VARCHAR(15) UNIQUE NOT NULL,
    email VARCHAR(30) UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
)
