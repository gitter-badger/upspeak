-- This migration restructures nodes data storage to ease frequent access
alter table public.nodes
    -- time when node last updated. Null value would mean it was never updated
    add column updated_at timestamp with time zone,
    -- author who last updated the node. Null value would mean it was never updated
    add column updated_by bigint references public.users(id),
    -- Subject of the node. Can be null based on data type
    add column subject text,
    -- Text body of the node. Can be null based on data type
    add column body text,
    -- Extra JSON data. Can be null based on data type. Used by data types that require rich data
    add column rich_data jsonb,
    -- Remove the revision_head column since we will have the latest data in the same table
    drop column revision_head cascade;

-- Store node_revisions
alter table public.node_revisions
    -- remove the revision ID because we can use timestamp and node ID to get what we need
    drop constraint node_revisions_pkey,
    drop column id,
    -- Rename extra to match new name for nodes
    -- rename column extra to rich_data,
    -- and set new primary key as combination of node_id and timestamp
    add primary key(node_id, created_at);

alter table public.node_revisions
    -- Rename extra to match new name for nodes
    rename column extra to rich_data;

-- Remove sequence for node_revisions
drop sequence public.node_revision_seq;
