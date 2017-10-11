CREATE TABLE IF NOT EXISTS messages (
	id SERIAL PRIMARY KEY,
	room_id INTEGER REFERENCES rooms,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

CREATE INDEX ON messages (room_id);
CREATE INDEX ON messages (created_at);
