\set ON_ERROR_STOP
\encoding utf8
SET client_min_messages TO INFO;

CREATE TABLE "user" (
    "id" serial PRIMARY KEY,
    "username" VARCHAR(255),
    CONSTRAINT "uq_username" UNIQUE("username"),
    "password" TEXT NOT NULL,
    registration_datetime TIMESTAMP(2) WITHOUT TIME ZONE NOT NULL DEFAULT current_timestamp(2)
);

CREATE TABLE user_session (
    session_key text primary key,
    user_id int not null REFERENCES "user"("id"), -- Could have a hard "references User"
    registration_datetime TIMESTAMP(2) WITHOUT TIME ZONE NOT NULL DEFAULT current_timestamp(2)
);

CREATE TABLE room(
    "id" serial PRIMARY KEY,
    "user_id" INT REFERENCES "user"("id") NOT NULL,
    "username" VARCHAR(255) REFERENCES "user"("username") NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    CONSTRAINT "uq_name" UNIQUE("username", "name")
);

CREATE TABLE "messages"(
    "id" serial PRIMARY KEY,
    "room_id" INT REFERENCES "room"("id") NOT NULL,
    "user_id" INT REFERENCES "user"("id") NOT NULL,
    "content" TEXT NOT NULL,
    "creation_datetime" TIMESTAMP(2) WITHOUT TIME ZONE NOT NULL DEFAULT current_timestamp(2),
    "edition_date" TIMESTAMP(2) WITHOUT TIME ZONE NOT NULL DEFAULT current_timestamp(2)
);

CREATE TABLE "badwords"(
    "id" serial PRIMARY KEY,
    "room_id" INT REFERENCES "room"("id") NOT NULL,
    "word" TEXT NOT NULL,
    "creation_datetime" TIMESTAMP(2) WITHOUT TIME ZONE NOT NULL DEFAULT current_timestamp(2)
);



INSERT INTO "message"("room_id", "user_id", "content")
    SELECT r.id, u.id, $4
        FROM room
        LEFT JOIN "user" u ON u.username = $3
        r WHERE r.name = $1 AND r.username = $2