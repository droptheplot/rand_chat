CREATE TABLE IF NOT EXISTS rooms (
	id serial,
	owner_id bigint NOT NULL,
	guest_id bigint,
	PRIMARY KEY(id)
);

CREATE UNIQUE INDEX ON rooms (owner_id);
