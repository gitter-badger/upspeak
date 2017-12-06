-- Enable plpgsql
create extension if not exists plpgsql with schema pg_catalog;

comment on extension plpgsql is 'PL/PGSQL procedural language';

-- We use the `public` schema by default
create schema if not exists public;

set search_path = public, pg_catalog;
