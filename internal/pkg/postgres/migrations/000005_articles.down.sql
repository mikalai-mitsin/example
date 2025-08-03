DROP TABLE public.articles;

DELETE
FROM public.permissions
WHERE id IN (
    'article_list',
    'article_detail',
    'article_create',
    'article_update',
    'article_delete'
);
