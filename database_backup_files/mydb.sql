-- -------------------------------------------------------------
-- TablePlus 6.6.8(632)
--
-- https://tableplus.com/
--
-- Database: mydb
-- Generation Time: 2025-08-25 19:40:18.2950
-- -------------------------------------------------------------


DROP TABLE IF EXISTS "public"."events";
-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS events_id_seq;

-- Table Definition
CREATE TABLE "public"."events" (
    "id" int4 NOT NULL DEFAULT nextval('events_id_seq'::regclass),
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."event_history";
-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS "eventHistory_id_seq";

-- Table Definition
CREATE TABLE "public"."event_history" (
    "id" int4 NOT NULL DEFAULT nextval('"eventHistory_id_seq"'::regclass),
    "player_id" int4 NOT NULL,
    "event_id" int4 NOT NULL,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."categories";
-- Table Definition
CREATE TABLE "public"."categories" (
    "id" int2 NOT NULL,
    "name" varchar NOT NULL,
    "min_level" int4,
    "max_level" int4,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."group_name";
-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS group_id_seq;

-- Table Definition
CREATE TABLE "public"."group_name" (
    "id" int4 NOT NULL DEFAULT nextval('group_id_seq'::regclass),
    "category_id" int2 NOT NULL,
    "event_id" int4 NOT NULL,
    "group_count" int4,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."players";
-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS players_id_seq;

-- Table Definition
CREATE TABLE "public"."players" (
    "id" int4 NOT NULL DEFAULT nextval('players_id_seq'::regclass),
    "event_number" int4,
    "email" varchar NOT NULL,
    "password_hash" varchar NOT NULL,
    "username" varchar NOT NULL,
    "level" int4 NOT NULL,
    "group" int4,
    PRIMARY KEY ("id")
);

INSERT INTO "public"."events" ("id") VALUES
(4);

INSERT INTO "public"."categories" ("id", "name", "min_level", "max_level") VALUES
(1, 'Bronze', 1, 20),
(2, 'Silver', 21, 49),
(3, 'Gold', 50, 150);

ALTER TABLE "public"."event_history" ADD FOREIGN KEY ("event_id") REFERENCES "public"."events"("id");
ALTER TABLE "public"."event_history" ADD FOREIGN KEY ("player_id") REFERENCES "public"."players"("id");


-- Indices
CREATE UNIQUE INDEX "eventHistory_pkey" ON public.event_history USING btree (id);


-- Indices
CREATE UNIQUE INDEX groups_pkey ON public.categories USING btree (id);
ALTER TABLE "public"."group_name" ADD FOREIGN KEY ("category_id") REFERENCES "public"."categories"("id");
ALTER TABLE "public"."group_name" ADD FOREIGN KEY ("event_id") REFERENCES "public"."events"("id");


-- Indices
CREATE UNIQUE INDEX group_pkey ON public.group_name USING btree (id);
ALTER TABLE "public"."players" ADD FOREIGN KEY ("group") REFERENCES "public"."group_name"("id");
ALTER TABLE "public"."players" ADD FOREIGN KEY ("event_number") REFERENCES "public"."events"("id");


-- Indices
CREATE UNIQUE INDEX players_username_key ON public.players USING btree (username);
CREATE UNIQUE INDEX players_email_key ON public.players USING btree (email);
