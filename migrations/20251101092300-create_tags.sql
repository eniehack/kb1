
-- +migrate Up
CREATE TABLE tags (
    id TEXT PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL
);

CREATE TABLE note_tag (
    note_id TEXT,
    tag_id TEXT,
    FOREIGN KEY (note_id) REFERENCES note(id),
    FOREIGN KEY (tag_id) REFERENCES tag(id)
);
-- +migrate Down
DROP TABLE note_tag;
DROP TABLE tags;