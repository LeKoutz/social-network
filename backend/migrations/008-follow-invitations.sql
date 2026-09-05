BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "invitations" (
    "id"    INTEGER NOT NULL UNIQUE,
    "timestamp"    TEXT NOT NULL,
    "status"    TEXT NOT NULL CHECK(status in ('pending', 'accepted', 'denied', 'unfollowed')) DEFAULT 'pending',
    "from_user_id"    INTEGER NOT NULL,
    "to_user_id"    INTEGER NOT NULL,
    PRIMARY KEY("id" AUTOINCREMENT),
    FOREIGN KEY("to_user_id") REFERENCES "users"("id") ON DELETE CASCADE,
    FOREIGN KEY("from_user_id") REFERENCES "users"("id") ON DELETE CASCADE,
    UNIQUE("from_user_id", "to_user_id")
);

END TRANSACTION;
