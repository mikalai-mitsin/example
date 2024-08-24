CREATE TABLE public.users
(
    id          uuid                  DEFAULT uuid_generate_v4()
        CONSTRAINT users_pk PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    password text NOT NULL,
    email text NOT NULL,
    group_id varchar NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);
CREATE INDEX search_users
    ON public.users
        USING GIN (to_tsvector('english', first_name || last_name || email));

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE
    ON
        public.users
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_task();

INSERT INTO public.permissions (id, name)
VALUES ('user_list', 'User list'),
       ('user_detail', 'User detail'),
       ('user_create', 'User create'),
       ('user_update', 'User update'),
       ('user_delete', 'User delete');

INSERT INTO public.group_permissions (group_id, permission_id)
VALUES ('admin', 'user_list'),
       ('admin', 'user_detail'),
       ('admin', 'user_create'),
       ('admin', 'user_update'),
       ('admin', 'user_delete'),
       ('user', 'user_list'),
       ('user', 'user_detail'),
       ('user', 'user_create'),
       ('user', 'user_update'),
       ('user', 'user_delete'),
       ('guest', 'user_list'),
       ('guest', 'user_detail');
