---
-- Run this file after running migrations to seed with dummy data
---

begin;

/**
 * Create dummy sequence
 */
create sequence if not exists dummy_insert cycle;

/**
 * Lorem Ipsum generator
 *
 * Source: https://stackoverflow.com/a/43743174/575242
 */
create or replace function lipsum( quantity_ integer ) returns character varying
    language plpgsql
    as $$
  declare
    words_       text[];
    returnValue_ text := '';
    random_      integer;
    ind_         integer;
  begin words_ := array['lorem', 'ipsum', 'dolor', 'sit', 'amet', 'consectetur',
  'adipiscing', 'elit', 'a', 'ac', 'accumsan', 'ad', 'aenean', 'aliquam',
  'aliquet', 'ante', 'aptent', 'arcu', 'at', 'auctor', 'augue', 'bibendum',
  'blandit', 'class', 'commodo', 'condimentum', 'congue', 'consequat',
  'conubia', 'convallis', 'cras', 'cubilia', 'cum', 'curabitur', 'curae',
  'cursus', 'dapibus', 'diam', 'dictum', 'dictumst', 'dignissim', 'dis',
  'donec', 'dui', 'duis', 'egestas', 'eget', 'eleifend', 'elementum', 'enim',
  'erat', 'eros', 'est', 'et', 'etiam', 'eu', 'euismod', 'facilisi',
  'facilisis', 'fames', 'faucibus', 'felis', 'fermentum', 'feugiat',
  'fringilla', 'fusce', 'gravida', 'habitant', 'habitasse', 'hac', 'hendrerit',
  'himenaeos', 'iaculis', 'id', 'imperdiet', 'in', 'inceptos', 'integer',
  'interdum', 'justo', 'lacinia', 'lacus', 'laoreet', 'lectus', 'leo', 'libero',
  'ligula', 'litora', 'lobortis', 'luctus', 'maecenas', 'magna', 'magnis',
  'malesuada', 'massa', 'mattis', 'mauris', 'metus', 'mi', 'molestie', 'mollis',
  'montes', 'morbi', 'mus', 'nam', 'nascetur', 'natoque', 'nec', 'neque',
  'netus', 'nibh', 'nisi', 'nisl', 'non', 'nostra', 'nulla', 'nullam', 'nunc',
  'odio', 'orci', 'ornare', 'parturient', 'pellentesque', 'penatibus', 'per',
  'pharetra', 'phasellus', 'placerat', 'platea', 'porta', 'porttitor',
  'posuere', 'potenti', 'praesent', 'pretium', 'primis', 'proin', 'pulvinar',
  'purus', 'quam', 'quis', 'quisque', 'rhoncus', 'ridiculus', 'risus', 'rutrum',
  'sagittis', 'sapien', 'scelerisque', 'sed', 'sem', 'semper', 'senectus',
  'sociis', 'sociosqu', 'sodales', 'sollicitudin', 'suscipit', 'suspendisse',
  'taciti', 'tellus', 'tempor', 'tempus', 'tincidunt', 'torquent', 'tortor',
  'tristique', 'turpis', 'ullamcorper', 'ultrices', 'ultricies', 'urna', 'ut',
  'varius', 'vehicula', 'vel', 'velit', 'venenatis', 'vestibulum', 'vitae',
  'vivamus', 'viverra', 'volutpat', 'vulputate'];
    for ind_ in 1 .. quantity_ loop
      ind_ := ( random() * ( array_upper( words_, 1 ) - 1 ) )::integer + 1;
      returnValue_ := returnValue_ || ' ' || words_[ind_];
    end loop;
    return returnValue_;
  end;
$$;

--
-- ----------------------------------------------
--

/**
 * Populate users table with 100,000 users
 */
insert into users
    (username, password, email_primary, created_at, is_verified, is_active, display_name)
    select
        lipsum(1) || i,
        '123456',
        lipsum(1) || i || '@upspeak.net',
        current_timestamp,
        true,
        true,
        initcap(lipsum(2))
    from generate_series(4, 100000) as i;

/**
 * Populate orgs table with 10,000 orgs
 */
insert into orgs
    (slug, name, primary_contact)
    select
        'company' ||  nextval('dummy_insert'),
        initcap(lipsum(2) || currval('dummy_insert') || ' Company '),
        id
    from users where random() < 0.01 limit 10000;

/**
 * Populate general team for each org - 10,000 default teams
 */
insert into teams
    (slug, display_name, org_id, permissions)
    select 'general', 'General', orgs.id, '{"team":"write", "org": "write", "public": "none"}'
    from orgs;

/**
 * Populate team_members table with 50,000 combinations
 */
insert into team_members (team_id, user_id)
    select teams.id, users.id
        from users, teams
    where random() < 0.01 limit 50000
    on conflict do nothing;

/**
 * Populate thread nodes - 100,000 threads
 */
do $$
begin
    for i in 1..100000 loop
        with ids as (
            -- Generate IDs before hand for nodes and node_revision
            select generate_id() as node_id,
                generate_id('node_revision_seq') as revision_id,

                -- These should be passed from application
                team_id, user_id
            from team_members where random() < 0.01 limit 1
        ),
        n as (
            -- Insert node first
            insert into nodes (id, author_id, data_type, revision_head, created_at)
                select node_id, user_id, 'markdown', revision_id, now() from ids
                returning created_at
        ),
        rev as (
            -- Insert node revision next
            insert into node_revisions(id, node_id, subject, body, committer_id)
                select revision_id, node_id, lipsum(2) || i || lipsum(3), lipsum(300), user_id from ids
        )
        insert into threads(id, team_id)
            select node_id, team_id from ids;
    end loop;
end;
$$;

drop sequence dummy_insert;
drop function lipsum(integer);

commit;
