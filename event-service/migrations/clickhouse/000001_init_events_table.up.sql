CREATE TABLE IF NOT EXISTS events
(
    event_id     UUID,
    event_type   LowCardinality(String),
    occurred_at  DateTime64(3),
    entity_id    Int64,
    payload      String  -- сырой JSON payload события
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(occurred_at)
ORDER BY (event_type, occurred_at)
TTL toDateTime(occurred_at) + INTERVAL 180 DAY;