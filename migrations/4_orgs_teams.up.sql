/**
 * Organizations info
 */
create table public.orgs (
    id bigint default generate_id() not null primary key,
    slug character varying(256) not null unique,
    primary_contact bigint not null references public.users (id),
    name text not null
);

/**
 * Teams are groups of users in organizations
 */
create table public.teams (
    id bigint default generate_id() not null primary key,
    org_id bigint not null references public.orgs (id),
    slug character varying(256) not null,
    display_name text,
    parent_team bigint references public.teams (id), -- Used to create hierarchy of teams. Keep null till implemented
    permissions jsonb not null,

    -- Combination of org_id and slug should be unique. Team slugs are unique in org scope
    unique(org_id, slug)
);

/**
 * Link table for storing lists of team members
 */
create table public.team_members (
    team_id bigint not null references public.teams(id),
    user_id bigint not null references public.users(id),
    access_level character varying(10), -- Can be 'read', 'write' or 'admin'

    -- Combo primary key
    primary key (team_id, user_id)
);

