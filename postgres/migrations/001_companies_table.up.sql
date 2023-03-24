CREATE TYPE company_type AS ENUM ('Corporations', 'NonProfit', 'Cooperative','Sole Proprietorship');

CREATE TABLE IF NOT EXISTS companies (
	id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(15) UNIQUE NOT NULL,
    description TEXT,
    amount_of_employees integer NOT NULL,
    registered boolean NOT NULL,
    type company_type NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
