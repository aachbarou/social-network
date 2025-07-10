CREATE TABLE IF NOT EXISTS "event" (
    "event_id" TEXT(255) not null,
    "group_id" TEXT(255) not null,
    "created_by" TEXT(255) not null,
    "created_at" DATETIME not null default CURRENT_TIMESTAMP,
    "title" TEXT(255) not null,
    "content" TEXT(255) not null,
    "date" DATETIME not null,
    primary key ("event_id")
);