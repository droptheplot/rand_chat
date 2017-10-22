DROP TYPE IF EXISTS app;
CREATE TYPE app AS ENUM ('telegram', 'vk');

CREATE TABLE IF NOT EXISTS rooms (
  id SERIAL PRIMARY KEY,
  owner_id BIGINT NOT NULL,
  owner_app app NOT NULL,
  guest_id BIGINT,
  guest_app app,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()

  CHECK (owner_id <> guest_id OR owner_app <> guest_app)
  CHECK ((guest_id IS NOT NULL AND guest_app IS NOT NULL) OR (guest_id IS NULL AND guest_app IS NULL))
);

CREATE UNIQUE INDEX ON rooms (owner_id) WHERE active = 'true';
