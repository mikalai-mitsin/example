DROP TABLE public.tags;

DELETE
FROM public.permissions
WHERE id IN (
    'tag_list',
    'tag_detail',
    'tag_create',
    'tag_update',
    'tag_delete'
);
