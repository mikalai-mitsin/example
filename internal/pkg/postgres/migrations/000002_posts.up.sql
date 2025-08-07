CREATE TABLE public.posts
(
    id          uuid                  DEFAULT uuidv7()
        CONSTRAINT posts_pk PRIMARY KEY,
    body text NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);
