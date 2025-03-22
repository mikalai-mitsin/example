DROP TABLE public.comments;

DELETE
FROM public.permissions
WHERE id IN (
    'comment_list',
    'comment_detail',
    'comment_create',
    'comment_update',
    'comment_delete'
);
