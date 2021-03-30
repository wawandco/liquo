CREATE TABLE IF NOT EXISTS public.databasechangelog (
	id             character varying(255)                  not null,
	author         character varying(255)                  not null,
	filename       character varying(255)                   not null,
	dateexecuted   timestamp without time zone             not null,
	orderexecuted  integer                                 not null,
	exectype       character varying(10)                   not null,
	md5sum         character varying(35),
	description    character varying(255),
	comments       character varying(255),
	tag            character varying(255),
	liquibase      character varying(20),
	contexts       character varying(255),
	labels         character varying(255),
	deployment_id  character varying(10)
);

CREATE TABLE IF NOT EXISTS public.databasechangeloglock (
	id           integer                                 not null,
	locked       boolean                                 not null,
	lockgranted  timestamp without time zone,
	lockedby     character varying(255)
);