package main

import (
	"github.com/droptheplot/rand_chat/env"
)

func main() {
	db, _ := env.Init()

	db.Exec(`
		DO $$
		DECLARE
			room_id integer;
			i integer = 0;
		BEGIN
			INSERT INTO rooms (owner_id, owner_app, guest_id, guest_app)
				VALUES (random() * 100, 'vk', random() * 100, 'telegram')
				RETURNING id INTO room_id;

			WHILE i < 10
			LOOP
				i := i + 1;

				INSERT INTO messages (room_id, created_at) VALUES
					(room_id, now() - (i * (random() * 1000) || ' minutes')::interval);
			END LOOP;
		END $$;
	`)
}
