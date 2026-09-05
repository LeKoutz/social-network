BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "posts_new" (
    "id"    INTEGER NOT NULL UNIQUE,
    "timestamp"    TEXT NOT NULL,
    "title"    TEXT NOT NULL,
    "body"    TEXT NOT NULL,
    "user_id"    INTEGER NOT NULL,
    "image_path" TEXT,
    "visibility" TEXT NOT NULL CHECK(visibility IN ('public','limited','private')) DEFAULT 'public',
    "group_id" INTEGER,
    PRIMARY KEY("id" AUTOINCREMENT),
    FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE,
    FOREIGN KEY("group_id") REFERENCES "groups"("id") ON DELETE CASCADE,
    CHECK (group_id IS NULL OR visibility = 'public')
);

INSERT INTO "posts_new" ("id", "timestamp", "title", "body", "user_id", "image_path")
SELECT "id", "timestamp", "title", "body", "user_id", "image_path"
FROM "posts";

DROP TABLE posts;

ALTER TABLE posts_new RENAME TO posts;

END TRANSACTION;
