-- Read-only structural signature of public; excludes row data and sequence counters.
SET search_path = pg_catalog, public;
WITH signature AS (
    SELECT jsonb_build_object(
        'relations', (SELECT jsonb_agg(jsonb_build_array(c.relname, c.relkind,
                c.relpersistence, c.relrowsecurity, c.relforcerowsecurity) ORDER BY c.relname)
            FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = 'public'),
        'columns', (SELECT jsonb_agg(jsonb_build_array(c.relname, a.attnum, a.attname,
                format_type(a.atttypid, a.atttypmod), a.attnotnull, a.attidentity, a.attgenerated,
                pg_get_expr(d.adbin, d.adrelid)) ORDER BY c.relname, a.attnum)
            FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
            JOIN pg_namespace n ON n.oid = c.relnamespace
            LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
            WHERE n.nspname = 'public' AND a.attnum > 0 AND NOT a.attisdropped),
        'constraints', (SELECT jsonb_agg(jsonb_build_array(c.conname, c.contype,
                pg_get_constraintdef(c.oid), c.convalidated) ORDER BY c.conname)
            FROM pg_constraint c JOIN pg_namespace n ON n.oid = c.connamespace WHERE n.nspname = 'public'),
        'indexes', (SELECT jsonb_agg(jsonb_build_array(c.relname, pg_get_indexdef(c.oid),
                i.indisvalid, i.indisready) ORDER BY c.relname)
            FROM pg_index i JOIN pg_class c ON c.oid = i.indexrelid
            JOIN pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = 'public'),
        'sequences', (SELECT jsonb_agg(jsonb_build_array(c.relname, s.seqtypid::regtype::text,
                s.seqstart, s.seqincrement, s.seqmax, s.seqmin, s.seqcache, s.seqcycle) ORDER BY c.relname)
            FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid
            JOIN pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = 'public'),
        'triggers', (SELECT jsonb_agg(pg_get_triggerdef(t.oid) ORDER BY c.relname, t.tgname)
            FROM pg_trigger t JOIN pg_class c ON c.oid = t.tgrelid
            JOIN pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = 'public' AND NOT t.tgisinternal),
        'routines', (SELECT count(*) FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace WHERE n.nspname = 'public'),
        'custom_types', (SELECT count(*) FROM pg_type t JOIN pg_namespace n ON n.oid = t.typnamespace
            WHERE n.nspname = 'public' AND t.typtype IN ('d', 'e')),
        'policies', (SELECT count(*) FROM pg_policy p JOIN pg_class c ON c.oid = p.polrelid
            JOIN pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = 'public')
    ) AS value
)
SELECT CASE
    WHEN to_regclass('public.flyway_schema_history') IS NOT NULL THEN 'MANAGED'
    WHEN value->'relations' = 'null'::jsonb AND value->>'routines' = '0'
        AND value->>'custom_types' = '0' THEN 'EMPTY'
    ELSE 'LEGACY:' || md5(value::text)
END FROM signature;
