/**
 * generate_id creates unique IDs based on timestamp
 *
 * Generates 64-bit integer (bigint) id using timestamp and a given sequence. 46
 * bits for timestamp, 5 bits for partition id (default 1), followed by 13 bits
 * of sequence id. This is a modification of Instagram''s Postgresql based id
 * generation system.
 */
create function public.generate_id(
        seq_name text default 'public.global_id_seq'::text,
        partition_id integer default 1, out result bigint
    ) returns bigint
    language plpgsql parallel safe
    as $$declare
    our_epoch bigint := 1512066600000;
    seq_id bigint;
    now_millis bigint;
begin
    select nextval(seq_name) % 8192 into seq_id;

    select floor(extract(epoch from clock_timestamp()) * 1000) into now_millis;
    result := (now_millis - our_epoch) << 18;
    result := result | (partition_id << 13);
    result := result | (seq_id);
end;
$$;

/**
 * global_id_seq is used to generate ids by generate_id() by default.
 */
create sequence public.global_id_seq
    start with 1
    increment by 1
    no minvalue
    no maxvalue
    cache 1
    cycle;
