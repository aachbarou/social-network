CREATE TABLE IF NOT EXISTS comments (
    "comment_id" TEXT(255) not null,
    "post_id" TEXT(255) not null,
    "created_at" datetime not null default CURRENT_TIMESTAMP,
    "created_by" TEXT(255) not null,
    "content" text null,
    "image" TEXT(255) null,
    primary key ("comment_id")
);
