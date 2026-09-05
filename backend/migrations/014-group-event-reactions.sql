BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "group_event_reactions" (
    "id"        INTEGER NOT NULL UNIQUE,
    "event_id"  INTEGER NOT NULL,
    "user_id"   INTEGER NOT NULL,
    "reaction"  TEXT NOT NULL CHECK(reaction IN ('attending', 'not_attending', 'maybe')),
    "timestamp" TEXT NOT NULL,
    PRIMARY KEY("id" AUTOINCREMENT),
    UNIQUE("event_id", "user_id"),
    FOREIGN KEY("event_id") REFERENCES "group_events"("id") ON DELETE CASCADE,
    FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);

END TRANSACTION;
