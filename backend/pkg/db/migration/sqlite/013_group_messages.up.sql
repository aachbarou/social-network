CREATE TABLE IF NOT EXISTS group_messages (
			"message_id" TEXT(255) not null,
			"receiver_id" TEXT(255) not null,
			"is_read" INT default 0
);

