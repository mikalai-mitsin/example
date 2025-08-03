DROP TABLE public.likes;

DELETE
FROM public.permissions
WHERE id IN (
    'like_list',
    'like_detail',
    'like_create',
    'like_update',
    'like_delete'
);
