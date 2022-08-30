CREATE TABLE IF NOT EXISTS employees(
    id UUID PRIMARY KEY,
    company_id UUID REFERENCES companies,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    salary BIGINT NOT NULL
);