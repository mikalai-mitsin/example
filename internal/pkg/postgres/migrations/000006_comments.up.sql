CREATE TABLE public.comments
(
    id          uuid                  DEFAULT uuid_generate_v4()
        CONSTRAINT comments_pk PRIMARY KEY,
    text text NOT NULL,
    author_id uuid NOT NULL,
    post_id uuid NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);
INSERT INTO public.permissions (id, name)
VALUES ('comment_list', 'Comment list'),
       ('comment_detail', 'Comment detail'),
       ('comment_create', 'Comment create'),
       ('comment_update', 'Comment update'),
       ('comment_delete', 'Comment delete');

INSERT INTO public.group_permissions (group_id, permission_id)
VALUES ('admin', 'comment_list'),
       ('admin', 'comment_detail'),
       ('admin', 'comment_create'),
       ('admin', 'comment_update'),
       ('admin', 'comment_delete'),
       ('user', 'comment_list'),
       ('user', 'comment_detail'),
       ('user', 'comment_create'),
       ('user', 'comment_update'),
       ('user', 'comment_delete'),
       ('guest', 'comment_list'),
       ('guest', 'comment_detail');