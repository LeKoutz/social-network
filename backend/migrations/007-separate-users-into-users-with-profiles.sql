BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "user_profiles" (
    "id" INTEGER NOT NULL UNIQUE,
    "user_id" INTEGER NOT NULL UNIQUE,
    "first_name" TEXT NOT NULL DEFAULT '',
    "last_name" TEXT NOT NULL DEFAULT '',
    "nickname" TEXT,
    "about" TEXT,
    "avatar_url" TEXT,
    "date_of_birth" TEXT,
    "private_profile" BOOLEAN NOT NULL DEFAULT true,
    PRIMARY KEY("id" AUTOINCREMENT),
    FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);

INSERT INTO "user_profiles" ("user_id", "first_name", "last_name")
SELECT "id", "first_name", "last_name"
FROM "users";

CREATE TABLE "users_new" (
    "id" INTEGER NOT NULL UNIQUE,
    "email" TEXT NOT NULL UNIQUE,
    "username" TEXT NOT NULL UNIQUE,
    "hash" TEXT,
    "session_key" TEXT,
    "oauth_provider" TEXT,
    PRIMARY KEY("id" AUTOINCREMENT)
);

INSERT INTO "users_new" ("id", "email", "username", "hash", "session_key", "oauth_provider")
SELECT "id", "email", "username", "hash", "session_key", "oauth_provider"
FROM "users";

DROP TABLE "users";
ALTER TABLE "users_new" RENAME TO "users";

END TRANSACTION;
