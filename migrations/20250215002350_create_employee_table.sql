-- +goose Up
-- +goose StatementBegin
CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255) NULL,
    state VARCHAR(255) NULL,
    designation VARCHAR(255),
    salary FLOAT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(50),
    hire_date TIMESTAMP NULL,
    date_of_birth TIMESTAMP NULL,
    gender VARCHAR(50),
    marital_status VARCHAR(50),
    children INT,
    emergency_contact VARCHAR(50),
    emergency_contact_relation VARCHAR(50),
    address TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE employees;
-- +goose StatementEnd
