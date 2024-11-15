DROP TABLE public.users;

DELETE
FROM public.permissions
WHERE id IN (
    'user_list',
    'user_detail',
    'user_create',
    'user_update',
    'user_delete'
);
