/**
 * Threads store thread-level states and are linked to teams. Thread ID is a
 * source-node ID.
 */
create table public.threads (
    id bigint not null primary key references public.nodes (id) deferrable initially deferred,
    -- team the thread belongs to. Same as source node's team ID
    team_id bigint not null references public.teams (id),
    permissions jsonb, -- thread permissions
     -- if null, thread is not forked
    forked_from_node bigint references public.nodes (id) deferrable initially deferred,
    merge_node bigint references public.nodes (id) deferrable initially deferred,
    is_open boolean default true, -- thread is open to new comments
    attrs jsonb -- thread-level attributes
);

-- Index teams to pull threads quicker
create index thread_team_idx on threads using btree (team_id);

-- gin indexing for thread attributes
create index thread_attrs_idx on threads using gin (attrs jsonb_path_ops);
