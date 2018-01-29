-- Adds archive capability (soft delete) to nodes and threads

alter table public.nodes
    add column is_archived bool not null default false;

alter table public.threads
    add column is_archived bool not null default false;
