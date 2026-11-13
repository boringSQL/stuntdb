\restrict tQ4nYb2sV1cLk8mR

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'customers',
	'relpages', '1820'::integer,
	'reltuples', '300000'::real,
	'relallvisible', '1810'::integer,
	'relallfrozen', '0'::integer
);

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'orders',
	'relpages', '41000'::integer,
	'reltuples', '2000000'::real,
	'relallvisible', '40900'::integer,
	'relallfrozen', '0'::integer
);

SELECT * FROM pg_catalog.pg_restore_relation_stats(
	'version', '180007'::integer,
	'schemaname', 'public',
	'relname', 'events',
	'relpages', '730'::integer,
	'reltuples', '96000'::real,
	'relallvisible', '700'::integer,
	'relallfrozen', '0'::integer
);

\unrestrict tQ4nYb2sV1cLk8mR
