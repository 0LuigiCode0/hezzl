-- +goose Up
CREATE TABLE IF NOT EXISTS log_goods_store
(
   id UInt64,
   project_id UInt64,
   name String,
   description String,
   priority Int64,  
   removed Bool,
   eventTime DateTime,
   
   index project_id_idx project_id type set(0),
   index name_idx name type bloom_filter
)
ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS log_goods_nats (
   id UInt64,
   project_id UInt64,
   name String,
   description String,
   priority Int64,
   removed Bool,
   eventTime DateTime
) ENGINE = NATS 
SETTINGS
   nats_url = 'nats://host.docker.internal:4222',
   nats_subjects = 'goods',
   nats_format = 'JSONEachRow',
   nats_username = 'admin',
   nats_password = 'admin',
   nats_max_block_size = 5000,
   nats_flush_interval_ms = 10000,
   nats_skip_broken_messages = 1000,
   nats_commit_mode = 'explicit';

CREATE MATERIALIZED VIEW IF NOT EXISTS log_goods_mv
TO log_goods_store AS
SELECT *
FROM log_goods_nats;

-- +goose Down
DROP TABLE IF EXISTS log_goods_nats;
DROP TABLE IF EXISTS log_goods_store;
DROP TABLE IF EXISTS log_goods_mv;