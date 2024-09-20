CREATE TABLE metrics (
    id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(255) NOT NULL DEFAULT '',
    delta BIGINT NOT NULL DEFAULT 0,
    value DOUBLE PRECISION NOT NULL DEFAULT 0
);