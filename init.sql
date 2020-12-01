BEGIN;

-- Table: public.room

-- DROP TABLE public.room;

CREATE TABLE public.room
(
    created timestamp without time zone,
    id text COLLATE pg_catalog."default" NOT NULL,
    variant_type integer,
    variant_name text,
    max_players integer,
    num_of_rounds integer,
    current_round integer,
    CONSTRAINT room_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.room
    OWNER to postgres;




-- Table: public.player

-- DROP TABLE public.player;

CREATE TABLE public.player
(
    created timestamp without time zone,
    player_id text COLLATE pg_catalog."default" NOT NULL,
    room_id text COLLATE pg_catalog."default" NOT NULL,
    in_room boolean,
    ready_status boolean,
    name text COLLATE pg_catalog."default",
    CONSTRAINT room_foreign_key FOREIGN KEY (room_id)
        REFERENCES public.room (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.player
    OWNER to postgres;


-- Table: public.bet

-- DROP TABLE public.bet;

CREATE TABLE public.bet
(
    created timestamp without time zone,
    round_no integer NOT NULL,
    bettype integer NOT NULL,
    stake numeric NOT NULL,
    liability numeric,
    player_id text COLLATE pg_catalog."default" NOT NULL,
    result integer,
    total_return integer,
    room_id text COLLATE pg_catalog."default" NOT NULL
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.bet
    OWNER to postgres;


-- Table: public.result

-- DROP TABLE public.result;

CREATE TABLE public.result
(
    created time without time zone,
    round_no integer NOT NULL,
    room_id text COLLATE pg_catalog."default" NOT NULL,
    "number" integer NOT NULL,
    colour integer NOT NULL
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.result
    OWNER to postgres;





COMMIT;