-- We use the `public` schema by default
create schema if not exists public;

set search_path = public, pg_catalog;
