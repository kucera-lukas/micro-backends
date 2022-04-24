CREATE TABLE "messages" (
    "id" SERIAL PRIMARY KEY NOT NULL,
    "data" character varying NOT NULL,
    "created" timestamp(0)with time zone NOT NULL,
    "modified" timestamp(0)with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_message_modified
BEFORE UPDATE ON messages
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
