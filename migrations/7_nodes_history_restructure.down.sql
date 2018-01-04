-- This migration is a breaking migration. We cannot recover all data, so the
-- `down` attempts to recreate the previous migration's state as much as
-- possible.

alter table public.nodes
    add column revision_head bigint, -- cannot add `not null` constraint here because no data is present now
    drop column updated_at cascade,
    drop column updated_by cascade,
    drop column subject cascade,
    drop column body cascade,
    drop column rich_data cascade;

alter table public.node_revisions
    -- remove the revision ID because we can use timestamp and node ID to get what we need
    drop constraint node_revisions_pkey;

-- bring back node_revisions_seq
create sequence public.node_revision_seq
    start with 1
    increment by 1
    no minvalue
    no maxvalue
    cache 1
    cycle;

-- Bring back id as primary key
alter table public.node_revisions
    add column id bigint default generate_id('public.node_revision_seq'::text) not null;

alter table public.node_revisions
    add primary key(id);

-- Bring back constraints for nodes.revision_head
alter table only public.nodes
    add constraint fk_node_revision_head foreign key (revision_head) references node_revisions(id)
    deferrable initially deferred;

alter table public.node_revisions
    -- Undo rename
    rename column rich_data to extra;
