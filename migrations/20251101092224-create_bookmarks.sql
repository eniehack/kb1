
-- +migrate Up
PRAGMA foreign_key=1;
CREATE TABLE notes (
    id TEXT PRIMARY KEY,
    content TEXT NOT NULL,
    url TEXT,
    is_public BOOLEAN DEFAULT false,
    created_at timestamp
);

CREATE TABLE read_laters (
    note_id TEXT NOT NULL,
    read BOOLEAN DEFAULT false,
    FOREIGN KEY (note_id) REFERENCES notes(id)
);
-- +migrate Down
DROP TABLE to_read;
DROP TABLE notes;