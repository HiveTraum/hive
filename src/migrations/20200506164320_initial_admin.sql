-- +goose Up
-- +goose StatementBegin
DO
$$
    DECLARE
        user_identifier UUID = uuid_generate_v4();
        role_identifier UUID = uuid_generate_v4();

    BEGIN
        --         insert into my_table (uuid_column) values (my_uuid);
--         select * from my_table where uuid_column = my_uuid;
        INSERT INTO users (id) VALUES (user_identifier);
        INSERT INTO roles (id, title) VALUES (role_identifier, 'admin');
        INSERT INTO emails (id, user_id, value) VALUES (uuid_generate_v4(), user_identifier, 'admin@admin.com');
        INSERT INTO user_roles (id, user_id, role_id) VALUES (uuid_generate_v4(), user_identifier, role_identifier);
        INSERT INTO passwords (id, user_id, value)
        VALUES (uuid_generate_v4(), user_identifier, '5a8fc31d-8a7c-4ffb-8d3b-20af8b619a0d');
    END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
