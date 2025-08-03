CREATE TABLE public.likes
(
    id          uuid                  DEFAULT uuid_generate_v4()
        CONSTRAINT likes_pk PRIMARY KEY,
    post_id uuid NOT NULL,
    value text NOT NULL,
    user_id uuid NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);
