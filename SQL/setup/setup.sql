-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler version: 1.1.1
-- PostgreSQL version: 16.0
-- Project Site: pgmodeler.io
-- Model Author: ---
-- Database creation must be performed outside a multi lined SQL file. 
-- These commands were put in this file only as a convenience.
-- 
-- object: "Knight-Link" | type: DATABASE --
-- DROP DATABASE IF EXISTS "Knight-Link";
-- CREATE DATABASE "Knight-Link";
-- ddl-end --
SET check_function_bodies = false;
-- ddl-end --
-- object: public.auth | type: TYPE --
-- DROP TYPE IF EXISTS public.auth CASCADE;
CREATE TYPE public.auth AS ENUM ('student', 'admin', 'superadmin');
-- ddl-end --
-- object: public."Universities" | type: TABLE --
-- DROP TABLE IF EXISTS public."Universities" CASCADE;
CREATE TABLE public."Universities" (
    uni_id serial NOT NULL,
    name varchar(255) NOT NULL,
    description text,
    student_no integer,
    picture bytea,
    CONSTRAINT "Universities_pk" PRIMARY KEY (uni_id),
    CONSTRAINT uni_ques UNIQUE (name)
);
-- ddl-end --
COMMENT ON COLUMN public."Universities".student_no IS E'Number of students in the university currently';
-- object: public."Locations" | type: TABLE --
-- DROP TABLE IF EXISTS public."Locations" CASCADE;
CREATE TABLE public."Locations" (
    loc_id serial NOT NULL,
    address text,
    latitude varchar(15) NOT NULL,
    longitude varchar(15) NOT NULL,
    CONSTRAINT "Locations_pk" PRIMARY KEY (loc_id)
);
-- ddl-end --
-- ddl-end --
-- object: public."Users" | type: TABLE --
-- DROP TABLE IF EXISTS public."Users" CASCADE;
CREATE TABLE public."Users" (
    user_id serial NOT NULL,
    first_name varchar(255),
    last_name varchar(255),
    email varchar(255),
    username varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    user_type public.auth NOT NULL,
    profile_picture bytea,
    uni_id serial,
    CONSTRAINT "Users_pk" PRIMARY KEY (user_id),
    CONSTRAINT unique_username UNIQUE (username)
);
-- ddl-end --
ALTER TABLE public."Users" ENABLE ROW LEVEL SECURITY;
-- ddl-end --
-- object: public.event | type: TYPE --
-- DROP TYPE IF EXISTS public.event CASCADE;
CREATE TYPE public.event AS ENUM ('public', 'private', 'rso_event');
-- ddl-end --
-- object: public."RSOs" | type: TABLE --
-- DROP TABLE IF EXISTS public."RSOs" CASCADE;
CREATE TABLE public."RSOs" (
    rso_id serial NOT NULL,
    name varchar(255) NOT NULL,
    description text,
    uni_id serial NOT NULL,
    admin_id serial NOT NULL,
    date_created timestamp with time zone NOT NULL,
    CONSTRAINT "RSOs_pk" PRIMARY KEY (rso_id),
    CONSTRAINT rso_uniques UNIQUE (name)
);
-- ddl-end --
ALTER TABLE public."RSOs" ENABLE ROW LEVEL SECURITY;
-- ddl-end --
-- object: public.categories | type: TYPE --
-- DROP TYPE IF EXISTS public.categories CASCADE;
CREATE TYPE public.categories AS ENUM ('social', 'fundraising', 'tech talk', 'academic');
-- ddl-end --
-- object: public."Events" | type: TABLE --
-- DROP TABLE IF EXISTS public."Events" CASCADE;
CREATE TABLE public."Events" (
    event_id serial NOT NULL,
    name varchar(255) NOT NULL,
    tags public.categories [],
    description text,
    start_time timestamp with time zone NOT NULL,
    end_time timestamp with time zone NOT NULL,
    loc_id serial,
    contact_phone varchar(15),
    contact_email varchar(255),
    visibility public.event NOT NULL,
    uni_id serial,
    rso_id serial,
    superadmin_approval boolean DEFAULT FALSE,
    CONSTRAINT "Events_pk" PRIMARY KEY (event_id)
);
-- ddl-end --
-- object: public."User_RSO_Membership" | type: TABLE --
-- DROP TABLE IF EXISTS public."User_RSO_Membership" CASCADE;
CREATE TABLE public."User_RSO_Membership" (
    user_id serial NOT NULL,
    rso_id serial NOT NULL,
    CONSTRAINT "User_RSO_Membership_pk" PRIMARY KEY (user_id, rso_id)
);
-- ddl-end --
-- object: public."RSO_Apps" | type: TABLE --
-- DROP TABLE IF EXISTS public."RSO_Apps" CASCADE;
CREATE TABLE public."RSO_Apps" (
    id serial NOT NULL,
    name varchar(255) NOT NULL,
    description text DEFAULT 'None given',
    uni_id serial NOT NULL,
    admin_id serial NOT NULL,
    student1_id serial NOT NULL,
    student2_id serial NOT NULL,
    student3_id serial NOT NULL,
    superadmin_approval boolean NOT NULL DEFAULT FALSE,
    CONSTRAINT no_dupes UNIQUE NULLS NOT DISTINCT (name),
    CONSTRAINT "RSO_Apps_pk" PRIMARY KEY (id)
);
-- ddl-end --
-- object: public."Event_Feedback" | type: TABLE --
-- DROP TABLE IF EXISTS public."Event_Feedback" CASCADE;
CREATE TABLE public."Event_Feedback" (
    fb_id serial NOT NULL,
    user_id serial NOT NULL,
    event_id serial NOT NULL,
    content text,
    rating smallint,
    feedback_type varchar(15) NOT NULL,
    "timestamp" timestamp DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT rating CHECK (
        rating >= 1
        AND rating <= 5
    ),
    CONSTRAINT fb_type CHECK (feedback_type IN ('comment', 'rating')),
    CONSTRAINT "Event_Feedback_pk" PRIMARY KEY (fb_id)
);
-- ddl-end --
-- object: public.validate_non_overlapping_events | type: FUNCTION --
-- DROP FUNCTION IF EXISTS public.validate_non_overlapping_events() CASCADE;
CREATE FUNCTION public.validate_non_overlapping_events() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN IF EXISTS (
    SELECT 1
    FROM public."Events" ev
    WHERE (
            -- Check for overlap with existing events
            (
                NEW.start_time < ev.end_time
                AND NEW.end_time > ev.start_time
            )
            OR (
                NEW.start_time <= ev.start_time
                AND NEW.end_time >= ev.end_time
            )
            OR (
                NEW.start_time >= ev.start_time
                AND NEW.end_time <= ev.end_time
            )
        )
        AND ev.event_id <> NEW.event_id -- Exclude the event being inserted/updated
) THEN RAISE EXCEPTION 'Event conflicts with existing event times';
END IF;
RETURN NEW;
END;
$$;
-- ddl-end --
-- object: validate_event_before_insert | type: TRIGGER --
-- DROP TRIGGER IF EXISTS validate_event_before_insert ON public."Events" CASCADE;
CREATE TRIGGER validate_event_before_insert BEFORE
INSERT ON public."Events" FOR EACH ROW EXECUTE PROCEDURE public.validate_non_overlapping_events();
-- ddl-end --
-- object: validate_before_update | type: TRIGGER --
-- DROP TRIGGER IF EXISTS validate_before_update ON public."Events" CASCADE;
CREATE TRIGGER validate_before_update BEFORE
UPDATE ON public."Events" FOR EACH STATEMENT EXECUTE PROCEDURE public.validate_non_overlapping_events();
-- ddl-end --
-- object: public.validate_admin_association | type: FUNCTION --
-- DROP FUNCTION IF EXISTS public.validate_admin_association() CASCADE;
CREATE FUNCTION public.validate_admin_association() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN IF NEW.user_type = 'admin' THEN IF NOT EXISTS (
    SELECT 1
    FROM public."User_RSO_Membership"
    WHERE user_id = NEW.user_id
) THEN RAISE EXCEPTION 'Admin must be associated with at least one RSO';
END IF;
END IF;
RETURN NEW;
END;
$$;
-- ddl-end --
-- object: validate_admin_before_insert | type: TRIGGER --
-- DROP TRIGGER IF EXISTS validate_admin_before_insert ON public."Users" CASCADE;
CREATE TRIGGER validate_admin_before_insert BEFORE
INSERT ON public."Users" FOR EACH ROW EXECUTE PROCEDURE public.validate_admin_association();
-- ddl-end --
-- object: validate_admin_before_update | type: TRIGGER --
-- DROP TRIGGER IF EXISTS validate_admin_before_update ON public."Users" CASCADE;
CREATE TRIGGER validate_admin_before_update BEFORE
UPDATE ON public."Users" FOR EACH ROW EXECUTE PROCEDURE public.validate_admin_association();
-- ddl-end --
-- object: public.get_student_count | type: FUNCTION --
-- DROP FUNCTION IF EXISTS public.get_student_count(serial) CASCADE;
CREATE OR REPLACE FUNCTION public.get_student_count (IN university_id integer) RETURNS integer LANGUAGE plpgsql VOLATILE CALLED ON NULL INPUT SECURITY INVOKER PARALLEL UNSAFE COST 1 AS $$
DECLARE count_result integer;
BEGIN
SELECT COUNT(*) INTO count_result
FROM public."Users"
WHERE uni_id = university_id;
RETURN count_result;
END;
$$;
-- ddl-end --
-- object: public.update_student_count | type: FUNCTION --
-- DROP FUNCTION IF EXISTS public.update_student_count() CASCADE;
CREATE OR REPLACE FUNCTION public.update_student_count () RETURNS trigger LANGUAGE plpgsql VOLATILE CALLED ON NULL INPUT SECURITY INVOKER PARALLEL UNSAFE COST 1 AS $$
DECLARE count_result integer;
BEGIN
SELECT public.get_student_count(NEW.uni_id) INTO count_result;
UPDATE public."Universities"
SET student_no = count_result
WHERE uni_id = NEW.uni_id;
RETURN NEW;
END;
$$;
-- ddl-end --
-- object: update_student_count | type: TRIGGER --
-- DROP TRIGGER IF EXISTS update_student_count ON public."Users" CASCADE;
CREATE TRIGGER update_student_count
AFTER
INSERT
    OR DELETE
    OR
UPDATE ON public."Users" FOR EACH ROW EXECUTE PROCEDURE public.update_student_count();
-- ddl-end --
-- object: uni_id | type: CONSTRAINT --
-- ALTER TABLE public."Users" DROP CONSTRAINT IF EXISTS uni_id CASCADE;
ALTER TABLE public."Users"
ADD CONSTRAINT uni_id FOREIGN KEY (uni_id) REFERENCES public."Universities" (uni_id) MATCH SIMPLE ON DELETE
SET NULL ON UPDATE CASCADE;
-- ddl-end --
-- object: uni_id | type: CONSTRAINT --
-- ALTER TABLE public."RSOs" DROP CONSTRAINT IF EXISTS uni_id CASCADE;
ALTER TABLE public."RSOs"
ADD CONSTRAINT uni_id FOREIGN KEY (uni_id) REFERENCES public."Universities" (uni_id) MATCH SIMPLE ON DELETE
SET NULL ON UPDATE CASCADE;
-- ddl-end --
-- object: owner_id | type: CONSTRAINT --
-- ALTER TABLE public."RSOs" DROP CONSTRAINT IF EXISTS owner_id CASCADE;
ALTER TABLE public."RSOs"
ADD CONSTRAINT owner_id FOREIGN KEY (admin_id) REFERENCES public."Users" (user_id) MATCH SIMPLE ON DELETE
SET NULL ON UPDATE CASCADE;
-- ddl-end --
-- object: uni | type: CONSTRAINT --
-- ALTER TABLE public."Events" DROP CONSTRAINT IF EXISTS uni CASCADE;
ALTER TABLE public."Events"
ADD CONSTRAINT uni FOREIGN KEY (uni_id) REFERENCES public."Universities" (uni_id) MATCH SIMPLE ON DELETE
SET NULL ON UPDATE CASCADE;
-- ddl-end --
-- object: rso | type: CONSTRAINT --
-- ALTER TABLE public."Events" DROP CONSTRAINT IF EXISTS rso CASCADE;
ALTER TABLE public."Events"
ADD CONSTRAINT rso FOREIGN KEY (rso_id) REFERENCES public."RSOs" (rso_id) MATCH SIMPLE ON DELETE
SET DEFAULT ON UPDATE CASCADE;
-- ddl-end --
-- object: loc | type: CONSTRAINT --
-- ALTER TABLE public."Events" DROP CONSTRAINT IF EXISTS loc CASCADE;
ALTER TABLE public."Events"
ADD CONSTRAINT loc FOREIGN KEY (loc_id) REFERENCES public."Locations" (loc_id) MATCH SIMPLE ON DELETE
SET NULL ON UPDATE CASCADE;
-- ddl-end --
-- object: "user" | type: CONSTRAINT --
-- ALTER TABLE public."User_RSO_Membership" DROP CONSTRAINT IF EXISTS "user" CASCADE;
ALTER TABLE public."User_RSO_Membership"
ADD CONSTRAINT "user" FOREIGN KEY (user_id) REFERENCES public."Users" (user_id) MATCH SIMPLE ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --
-- object: rso | type: CONSTRAINT --
-- ALTER TABLE public."User_RSO_Membership" DROP CONSTRAINT IF EXISTS rso CASCADE;
ALTER TABLE public."User_RSO_Membership"
ADD CONSTRAINT rso FOREIGN KEY (rso_id) REFERENCES public."RSOs" (rso_id) MATCH SIMPLE ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --
-- object: uni | type: CONSTRAINT --
-- ALTER TABLE public."RSO_Apps" DROP CONSTRAINT IF EXISTS uni CASCADE;
ALTER TABLE public."RSO_Apps"
ADD CONSTRAINT uni FOREIGN KEY (uni_id) REFERENCES public."Universities" (uni_id) MATCH SIMPLE ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --
-- object: admin | type: CONSTRAINT --
-- ALTER TABLE public."RSO_Apps" DROP CONSTRAINT IF EXISTS admin CASCADE;
ALTER TABLE public."RSO_Apps"
ADD CONSTRAINT admin FOREIGN KEY (admin_id) REFERENCES public."Users" (user_id) MATCH SIMPLE ON DELETE
SET NULL ON UPDATE CASCADE;
-- ddl-end --
-- object: s1 | type: CONSTRAINT --
-- ALTER TABLE public."RSO_Apps" DROP CONSTRAINT IF EXISTS s1 CASCADE;
ALTER TABLE public."RSO_Apps"
ADD CONSTRAINT s1 FOREIGN KEY (student1_id) REFERENCES public."Users" (user_id) MATCH SIMPLE ON DELETE
SET NULL ON UPDATE CASCADE;
-- ddl-end --
-- object: s2 | type: CONSTRAINT --
-- ALTER TABLE public."RSO_Apps" DROP CONSTRAINT IF EXISTS s2 CASCADE;
ALTER TABLE public."RSO_Apps"
ADD CONSTRAINT s2 FOREIGN KEY (student2_id) REFERENCES public."Users" (user_id) MATCH SIMPLE ON DELETE
SET NULL ON UPDATE CASCADE;
-- ddl-end --
-- object: s3 | type: CONSTRAINT --
-- ALTER TABLE public."RSO_Apps" DROP CONSTRAINT IF EXISTS s3 CASCADE;
ALTER TABLE public."RSO_Apps"
ADD CONSTRAINT s3 FOREIGN KEY (student3_id) REFERENCES public."Users" (user_id) MATCH SIMPLE ON DELETE
SET NULL ON UPDATE CASCADE;
-- ddl-end --
-- object: "user" | type: CONSTRAINT --
-- ALTER TABLE public."Event_Feedback" DROP CONSTRAINT IF EXISTS "user" CASCADE;
ALTER TABLE public."Event_Feedback"
ADD CONSTRAINT "user" FOREIGN KEY (user_id) REFERENCES public."Users" (user_id) MATCH SIMPLE ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --
-- object: event | type: CONSTRAINT --
-- ALTER TABLE public."Event_Feedback" DROP CONSTRAINT IF EXISTS event CASCADE;
ALTER TABLE public."Event_Feedback"
ADD CONSTRAINT event FOREIGN KEY (event_id) REFERENCES public."Events" (event_id) MATCH SIMPLE ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --
-- Add Online Location for Usage
INSERT INTO public."Locations" (address, latitude, longitude)
VALUES ('Online', 'o', 'o');