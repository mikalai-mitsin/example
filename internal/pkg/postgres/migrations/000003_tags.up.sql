CREATE TABLE public.tags
(
    id          uuid                  DEFAULT uuidv7()
        CONSTRAINT tags_pk PRIMARY KEY,
    post_id uuid NOT NULL,
    value text NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);
