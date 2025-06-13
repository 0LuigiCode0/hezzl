-- +goose Up
CREATE TABLE IF NOT EXISTS projects (
   id bigserial PRIMARY KEY NOT NULL,
   name character varying(255) UNIQUE NOT NULL,
   created_at timestamp DEFAULT now() NOT NULL
);

INSERT INTO projects (name) VALUES ('Первая запись');

CREATE SEQUENCE IF NOT EXISTS goods_priority_seq MINVALUE 1;

CREATE TABLE IF NOT EXISTS goods (
   id bigserial PRIMARY KEY NOT NULL,
   project_id bigserial REFERENCES projects(id) NOT NULL,
   name character varying(255) UNIQUE NOT NULL,
   description text DEFAULT '' NOT NULL,
   priority integer DEFAULT nextval('goods_priority_seq') NOT NULL,
   removed boolean DEFAULT false NOT NULL,
   created_at timestamp DEFAULT now() NOT NULL
);

ALTER SEQUENCE IF EXISTS goods_priority_seq OWNED BY goods.priority;

-- +goose Down
DROP TABLE IF EXISTS goods;
DROP TABLE IF EXISTS projects;