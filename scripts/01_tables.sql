-- Drop table

-- DROP TABLE public.arquivos;

CREATE TABLE public.arquivos (
	id serial4 NOT NULL,
	file_md5 varchar(32) NOT NULL,
	file_path text NOT NULL,
	file_name text NOT NULL,
	file_entry text NOT NULL,
	ext text NULL,
	ext_U text NULL,
	CONSTRAINT arquivos_pk PRIMARY KEY (id),
	CONSTRAINT arquivos_un UNIQUE (file_md5, file_path)
);
CREATE INDEX arquivos_file_md5_idx ON public.arquivos USING btree (file_md5);