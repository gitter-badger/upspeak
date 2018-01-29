alter table public.nodes
    drop column is_archived;

alter table public.threads
    drop column is_archived;
