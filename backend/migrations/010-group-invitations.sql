BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "group_invitations" (
    "id"            INTEGER NOT NULL UNIQUE,
    "timestamp"     TEXT NOT NULL,
    "status"        TEXT NOT NULL CHECK(status IN ('pending', 'accepted', 'denied', 'left', 'removed')),
    "group_id"      INTEGER NOT NULL,
    "to_user_id"    INTEGER NOT NULL,
    "invited_by"    INTEGER,
    PRIMARY KEY("id" AUTOINCREMENT),
    UNIQUE("group_id", "to_user_id"),
    FOREIGN KEY("group_id") REFERENCES "groups"("id") ON DELETE CASCADE,
    FOREIGN KEY("to_user_id") REFERENCES "users"("id") ON DELETE CASCADE,
    FOREIGN KEY("invited_by") REFERENCES "users"("id") ON DELETE SET NULL
);

END TRANSACTION;
