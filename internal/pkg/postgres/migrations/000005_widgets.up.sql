CREATE TABLE public.widgets
(
    id          uuid                  DEFAULT uuid_generate_v4()
        CONSTRAINT widgets_pk PRIMARY KEY,
    form_screen_id text NOT NULL,
    name varchar NOT NULL,
    ordering bigint NOT NULL,
    is_optional boolean NOT NULL,
    ui_settings text NOT NULL,
    deleted_at timestamp NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE TRIGGER update_widgets_updated_at
    BEFORE UPDATE
    ON
        public.widgets
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_task();

INSERT INTO public.permissions (id, name)
VALUES ('widget_list', 'Widget list'),
       ('widget_detail', 'Widget detail'),
       ('widget_create', 'Widget create'),
       ('widget_update', 'Widget update'),
       ('widget_delete', 'Widget delete');

INSERT INTO public.group_permissions (group_id, permission_id)
VALUES ('admin', 'widget_list'),
       ('admin', 'widget_detail'),
       ('admin', 'widget_create'),
       ('admin', 'widget_update'),
       ('admin', 'widget_delete'),
       ('user', 'widget_list'),
       ('user', 'widget_detail'),
       ('user', 'widget_create'),
       ('user', 'widget_update'),
       ('user', 'widget_delete'),
       ('guest', 'widget_list'),
       ('guest', 'widget_detail');
