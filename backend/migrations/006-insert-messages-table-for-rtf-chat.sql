CREATE TABLE IF NOT EXISTS "messages" (
	"id"	INTEGER NOT NULL UNIQUE,
    "sender_id"	INTEGER NOT NULL,
    "recipient_id" INTEGER NOT NULL,
    "body" TEXT NOT NULL,
	"timestamp"	TEXT NOT NULL,
    "read" BOOLEAN NOT NULL DEFAULT 0,
	PRIMARY KEY("id" AUTOINCREMENT)
    FOREIGN KEY("sender_id") REFERENCES "users"("id")
    FOREIGN KEY("recipient_id") REFERENCES "users"("id")
);