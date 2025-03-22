DROP TABLE public.posts;

DELETE
FROM public.permissions
WHERE id IN (
    'post_list',
    'post_detail',
    'post_create',
    'post_update',
    'post_delete'
);
