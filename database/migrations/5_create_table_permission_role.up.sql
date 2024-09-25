CREATE TABLE permission_role (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    permission_id UUID NOT NULL,
    role_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (permission_id) REFERENCES permission (id),
    FOREIGN KEY (role_id) REFERENCES role (id)
);