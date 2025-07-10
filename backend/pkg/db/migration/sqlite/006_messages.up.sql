CREATE TABLE IF NOT EXISTS messages (
			"message_id" TEXT(255) not null,
			"sender_id" TEXT(255) not null,
			"receiver_id" TEXT(255) not null,
			"type" TEXT(255) not null,
			"created_at" datetime not null default CURRENT_TIMESTAMP,
			"content" TEXT(255) not null,
			"is_read" INT default 0,
			primary key ("message_id")
		);
