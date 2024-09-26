CREATE TABLE IF NOT EXISTS permission(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    method VARCHAR(255) NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    module_id UUID NOT NULL,
    FOREIGN KEY (module_id) REFERENCES module (id)
);
