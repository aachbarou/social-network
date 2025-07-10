CREATE TABLE IF NOT EXISTS users (
    "user_id" TEXT not null,
    "created_at" datetime not null default CURRENT_TIMESTAMP,
    "email" TEXT not null,
    "first_name" TEXT not null,
    "last_name" TEXT not null,
    "nickname" TEXT(255) null,
    "birthday" datetime not null,
    "image" TEXT(255) null,
    "about" TEXT null,
    "status" TEXT(255) not null default PUBLIC,
    "password" TEXT(255) not null,
    primary key ("user_id")
);
