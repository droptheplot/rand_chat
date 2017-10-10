CREATE TABLE IF NOT EXISTS rooms (
	id SERIAL PRIMARY KEY,
	owner_id BIGINT NOT NULL,
	guest_id BIGINT,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
  CHECK (owner_id <> guest_id)
);

CREATE UNIQUE INDEX ON rooms (owner_id) WHERE active = 'true';
