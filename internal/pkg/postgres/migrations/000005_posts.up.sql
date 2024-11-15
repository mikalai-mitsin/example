CREATE TABLE public.posts
(
    id          uuid                  DEFAULT uuid_generate_v4()
        CONSTRAINT posts_pk PRIMARY KEY,
    title varchar NOT NULL,
    order bigint NOT NULL,
    is_optional boolean NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE TRIGGER update_posts_updated_at
    BEFORE UPDATE
    ON
        public.posts
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_task();

INSERT INTO public.permissions (id, name)
VALUES ('post_list', 'Post list'),
       ('post_detail', 'Post detail'),
       ('post_create', 'Post create'),
       ('post_update', 'Post update'),
       ('post_delete', 'Post delete');

INSERT INTO public.group_permissions (group_id, permission_id)
VALUES ('admin', 'post_list'),
       ('admin', 'post_detail'),
       ('admin', 'post_create'),
       ('admin', 'post_update'),
       ('admin', 'post_delete'),
       ('user', 'post_list'),
       ('user', 'post_detail'),
       ('user', 'post_create'),
       ('user', 'post_update'),
       ('user', 'post_delete'),
       ('guest', 'post_list'),
       ('guest', 'post_detail');
