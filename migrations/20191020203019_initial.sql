-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id      UUID   DEFAULT uuid_generate_v4() PRIMARY KEY,
    created bigint DEFAULT extract(epoch from now()) * 1000
);

CREATE TABLE passwords
(
    id      UUID   DEFAULT uuid_generate_v4() PRIMARY KEY,
    created BIGINT DEFAULT extract(epoch from now()) * 1000,
    user_id UUID,
    value   varchar(255),
    FOREIGN KEY (user_id) REFERENCES users
);

CREATE TABLE emails
(
    id      UUID   DEFAULT uuid_generate_v4() PRIMARY KEY,
    created BIGINT DEFAULT extract(epoch from now()) * 1000,
    user_id UUID,
    value   varchar(255) UNIQUE,
    FOREIGN KEY (user_id) REFERENCES users
);

CREATE TABLE phones
(
    id      UUID   DEFAULT uuid_generate_v4() PRIMARY KEY,
    created BIGINT DEFAULT extract(epoch from now()) * 1000,
    user_id UUID,
    value   varchar(50) UNIQUE,
    FOREIGN KEY (user_id) REFERENCES users
);

CREATE TABLE roles
(
    id      UUID   DEFAULT uuid_generate_v4() PRIMARY KEY,
    created BIGINT DEFAULT extract(epoch from now()) * 1000,
    title   VARCHAR(255) UNIQUE
);

CREATE TABLE user_roles
(
    id      UUID   DEFAULT uuid_generate_v4() PRIMARY KEY,
    created BIGINT DEFAULT extract(epoch from now()) * 1000,
    user_id UUID,
    role_id UUID,
    FOREIGN KEY (user_id) REFERENCES users,
    FOREIGN KEY (role_id) REFERENCES roles
);

CREATE TABLE users_view
(
    id      UUID UNIQUE,
    created BIGINT,
    updated BIGINT,
    phones  TEXT[],
    roles   TEXT[],
    emails  TEXT[],
    role_id BIGINT[]
);

CREATE TABLE secrets
(
    id      UUID   DEFAULT uuid_generate_v4() PRIMARY KEY,
    created BIGINT DEFAULT extract(epoch from now()) * 1000,
    value   UUID   DEFAULT uuid_generate_v4()
);

CREATE TABLE sessions
(
    refresh_token UUID   DEFAULT uuid_generate_v4(),
    fingerprint   VARCHAR(200),
    user_id       UUID,
    secret_id     UUID,
    created       BIGINT DEFAULT extract(epoch from now()) * 1000,
    user_agent    TEXT,
    FOREIGN KEY (user_id) REFERENCES users,
    FOREIGN KEY (secret_id) REFERENCES secrets
);

CREATE INDEX ON phones (user_id);
CREATE INDEX ON emails (user_id);
CREATE INDEX ON user_roles (role_id);
CREATE INDEX ON user_roles (user_id);
CREATE INDEX ON passwords (user_id);
CREATE INDEX ON sessions (user_id);
CREATE INDEX ON sessions (secret_id);

CREATE INDEX user_views_role_id_idx on users_view USING GIN (role_id);
CREATE INDEX user_views_phones_idx on users_view USING GIN (phones);
CREATE INDEX sessions_fingerprint_idx on sessions (fingerprint);
CREATE INDEX sessions_refresh_token_idx on sessions (refresh_token);
CREATE UNIQUE INDEX user_roles_idx ON user_roles (user_id, role_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE emails;
DROP TABLE passwords;
DROP TABLE phones;
DROP TABLE users;
DROP TABLE roles;
DROP TABLE user_roles;
-- +goose StatementEnd
