CREATE TABLE public.articles
(
    id          uuid                  DEFAULT uuidv7()
        CONSTRAINT articles_pk PRIMARY KEY,
    title text NOT NULL,
    subtitle text NOT NULL,
    body text NOT NULL,
    is_published text NOT NULL,
    updated_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at  timestamp    NOT NULL DEFAULT (now() at time zone 'utc')
);
