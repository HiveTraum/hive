-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id      BIGSERIAL PRIMARY KEY,
    created bigint DEFAULT extract(epoch from now()) * 1000
);

CREATE TABLE passwords
(
    id      BIGSERIAL PRIMARY KEY,
    created bigint DEFAULT extract(epoch from now()) * 1000,
    user_id bigint,
    value   varchar(255),
    FOREIGN KEY (user_id) REFERENCES users
);

CREATE TABLE emails
(
    id      BIGSERIAL PRIMARY KEY,
    created BIGINT DEFAULT extract(epoch from now()) * 1000,
    user_id bigint,
    value   varchar(255) UNIQUE,
    FOREIGN KEY (user_id) REFERENCES users
);

CREATE TABLE phones
(
    id      BIGSERIAL PRIMARY KEY,
    created BIGINT DEFAULT extract(epoch from now()) * 1000,
    user_id BIGINT,
    value   varchar(50) UNIQUE,
    FOREIGN KEY (user_id) REFERENCES users
);

CREATE TABLE roles
(
    id      BIGSERIAL PRIMARY KEY,
    created BIGINT DEFAULT extract(epoch from now()) * 1000,
    title   VARCHAR(255) UNIQUE
);

CREATE TABLE user_roles
(
    id      BIGSERIAL PRIMARY KEY,
    created BIGINT DEFAULT extract(epoch from now()) * 1000,
    user_id BIGINT,
    role_id BIGINT,
    FOREIGN KEY (user_id) REFERENCES users,
    FOREIGN KEY (role_id) REFERENCES roles
);

CREATE TABLE users_view
(
    id      BIGINT UNIQUE,
    created BIGINT,
    updated BIGINT,
    phones  TEXT[],
    roles   TEXT[],
    emails  TEXT[]
);
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
