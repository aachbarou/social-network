CREATE TABLE IF NOT EXISTS "sessions" (
    "session_id" TEXT(255) NOT NULL PRIMARY KEY,
    "user_id" TEXT(255) NOT NULL,
    "expiration_time" DATETIME NOT NULL
);

