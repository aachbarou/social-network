CREATE TABLE IF NOT EXISTS posts (
    "post_id" TEXT(255) not null,
    "group_id" TEXT(255) null,
    "created_by" TEXT(255) not null,
    "created_at" datetime not null default CURRENT_TIMESTAMP,
    "content" TEXT null,
    "image" TEXT(255) null,
    "visibility" TEXT(255) null default PUBLIC,
    primary key ("post_id")
);
