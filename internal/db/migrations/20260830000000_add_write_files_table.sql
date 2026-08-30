-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS write_files (
    session_id TEXT NOT NULL CHECK (session_id != ''),
    path TEXT NOT NULL CHECK (path != ''),
    written_at INTEGER NOT NULL,  -- Unix timestamp in seconds when file was last written
    FOREIGN KEY (session_id) REFERENCES sessions (id) ON DELETE CASCADE,
    PRIMARY KEY (path, session_id)
);

CREATE INDEX IF NOT EXISTS idx_write_files_session_id ON write_files (session_id);
CREATE INDEX IF NOT EXISTS idx_write_files_path ON write_files (path);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_write_files_path;
DROP INDEX IF EXISTS idx_write_files_session_id;
DROP TABLE IF EXISTS write_files;
-- +goose StatementEnd
