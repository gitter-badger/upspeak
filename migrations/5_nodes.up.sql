/**
 * List of available node types
 */
create table public.node_data_types (
    id character varying(256) not null primary key,
    is_editable boolean default true,

    -- Optional attribute storage for node data type. This can be used to store
    -- data type specific information. This field is not indexed.
    attrs jsonb
);

-- create a `markdown` node type by default
insert into public.node_data_types
    values ('markdown', true, null);

/**
 * Create nodes table. Our central data structure
 */
create table nodes (
    id bigint default generate_id() not null primary key,
    author_id bigint not null references public.users (id),
    team_id bigint references public.teams (id),
    data_type character varying(256) not null references public.node_data_types(id),
    source_node_id bigint references public.nodes (id),
    created_at timestamp with time zone not null,
    in_reply_to bigint references public.nodes(id),
    attrs jsonb,
    revision_head bigint not null
);

-- Btree index of authors of nodes
create index node_author_idx on public.nodes using btree (author_id);

-- Btree index of teams that nodes belong to
create index node_team_idx on public.nodes using btree (team_id);

-- Btree index of source nodes to speed up thread lookups
create index node_source_node_idx on public.nodes using btree (source_node_id);

-- GIN indexing for node attrs
create index node_attrs_idx on public.nodes using gin (attrs jsonb_path_ops);

/**
 * Create special sequence for node revisions
 */
create sequence public.node_revision_seq
    start with 1
    increment by 1
    no minvalue
    no maxvalue
    cache 1
    cycle;

/**
 * Create node revisions table
 */
create table public.node_revisions (
    id bigint default generate_id('public.node_revision_seq'::text) not null primary key,
    node_id bigint not null,
    subject text,
    body text,
    extra jsonb,
    created_at timestamp with time zone default now() not null,
    committer_id bigint not null references public.users(id)
);

-- Btree index of node IDs in node_revisions to quickly lookup all revisions of a given node
create index node_revisions_node_idx on public.node_revisions using btree (node_id);

-- Create GIN index on node extra
create index node_revisions_extra_idx on public.node_revisions using gin (extra jsonb_path_ops);

-- Link nodes -> node_revisions using foreign keys
alter table only public.nodes
    add constraint fk_node_revision_head foreign key (revision_head) references node_revisions(id)
    deferrable initially deferred;

-- Link node_revisions -> nodes using foreign key
alter table only public.node_revisions
    add constraint fk_node_revision_node foreign key (node_id) references nodes(id)
    deferrable initially deferred;

