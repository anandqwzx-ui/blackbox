CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT REFERENCES users(id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT REFERENCES users(id),
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT REFERENCES users(id)
);

CREATE UNIQUE INDEX users_uuid_idx ON users (uuid);

CREATE TABLE organisations (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    owner_user_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT REFERENCES users(id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT REFERENCES users(id),
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT REFERENCES users(id)
);

CREATE UNIQUE INDEX organisations_uuid_idx ON organisations (uuid);
CREATE INDEX organisations_owner_idx ON organisations (owner_user_id) WHERE deleted_at IS NULL;

CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    organisation_id BIGINT NOT NULL REFERENCES organisations(id),
    owner_user_id BIGINT NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    description TEXT,
    code VARCHAR(5) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT REFERENCES users(id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT REFERENCES users(id),
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT REFERENCES users(id)
);

ALTER TABLE projects ADD CONSTRAINT projects_code_unique UNIQUE (code);
CREATE UNIQUE INDEX projects_uuid_idx ON projects (uuid);
CREATE INDEX projects_owner_idx ON projects (owner_user_id) WHERE deleted_at IS NULL;

CREATE TABLE endpoints (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    project_id BIGINT NOT NULL REFERENCES projects(id),
    name TEXT,
    method VARCHAR(16) NOT NULL,
    path TEXT NOT NULL,
    response_status INTEGER NOT NULL,
    response_body TEXT NOT NULL DEFAULT '',
    response_headers JSONB NOT NULL DEFAULT '{}'::jsonb,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT REFERENCES users(id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT REFERENCES users(id),
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT REFERENCES users(id)
);

CREATE UNIQUE INDEX endpoints_uuid_idx ON endpoints (uuid);
CREATE UNIQUE INDEX endpoints_project_path_method_unique ON endpoints (project_id, method, path) WHERE deleted_at IS NULL;
