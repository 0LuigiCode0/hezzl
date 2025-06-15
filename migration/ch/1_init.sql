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

-- +goose Down
DROP TABLE IF EXISTS log_goods;
