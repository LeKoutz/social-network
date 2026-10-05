BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "group_messages" (
    "id"        INTEGER NOT NULL UNIQUE,
    "group_id"  INTEGER NOT NULL,
    "sender_id" INTEGER NOT NULL,
    "body"      TEXT NOT NULL,
    "timestamp" TEXT NOT NULL,
    PRIMARY KEY("id" AUTOINCREMENT),
    FOREIGN KEY("group_id") REFERENCES "groups"("id") ON DELETE CASCADE,
    FOREIGN KEY("sender_id") REFERENCES "users"("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "group_message_reads" (
    "group_id"             INTEGER NOT NULL,
    "user_id"              INTEGER NOT NULL,
    "last_read_message_id" INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY("group_id", "user_id"),
    FOREIGN KEY("group_id") REFERENCES "groups"("id") ON DELETE CASCADE,
    FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);

END TRANSACTION;