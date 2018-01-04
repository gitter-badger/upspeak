-- Trigger function to keep node_revisions updated with node history

create or replace function trigger_on_node_revision()
    returns trigger
    language plpgsql as $body$
begin
    if tg_op = 'UPDATE' then
        if old.subject <> new.subject or old.body <> new.body or old.rich_data <> new.rich_data then
            -- Create node revision only if node's content has changed
            if old.updated_at is null and old.updated_by is null then
                -- First edit of node
                insert into audit.node_revisions (node_id, subject, body, rich_data, created_at, committer_id)
                values (old.id, old.subject, old.body, old.rich_data, old.created_at, old.author_id);
            else
                -- Subsequent edits of node
                insert into audit.node_revisions (node_id, subject, body, rich_data, created_at, committer_id)
                values (old.id, old.subject, old.body, old.rich_data, old.updated_at, old.updated_by);
            end if;
        end if;
        return new;
    end if;

    if tg_op = 'DELETE' then
        delete from audit.node_revisions where node_id = old.id;
        return old;
    end if;
end; $body$;

create trigger trigger_node_revision
  before update or delete
  on public.nodes
  for each row
  execute procedure trigger_on_node_revision();
