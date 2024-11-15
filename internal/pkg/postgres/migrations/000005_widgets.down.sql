DROP TABLE public.widgets;

DELETE
FROM public.permissions
WHERE id IN (
    'widget_list',
    'widget_detail',
    'widget_create',
    'widget_update',
    'widget_delete'
);
