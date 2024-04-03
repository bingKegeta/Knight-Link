CREATE TYPE user_type AS ENUM ('superadmin', 'admin', 'student');
CREATE TYPE event_type AS ENUM ('public', 'private', 'rso_event');
CREATE TABLE Universities (
    university_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    location VARCHAR(255) NOT NULL,
    description TEXT,
    number_of_students INTEGER,
    picture BYTEA
);
CREATE TABLE Locations (
    location_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(10, 8) NOT NULL
);
CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    auth user_type NOT NULL DEFAULT 'student',
    is_affiliated_with_rso BOOLEAN DEFAULT FALSE -- Optional flag for student affiliation
);
CREATE TABLE RSOs (
    rso_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    university_id INTEGER REFERENCES Universities(university_id) NOT NULL,
    admin_id INTEGER REFERENCES Users(user_id) NOT NULL
);
CREATE TABLE User_RSO_Membership (
    user_id INTEGER REFERENCES Users(user_id) NOT NULL,
    rso_id INTEGER REFERENCES RSOs(rso_id) NOT NULL,
    PRIMARY KEY (user_id, rso_id) -- Define composite primary key for both user and RSO
);
CREATE TABLE Events (
    event_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    event_category VARCHAR(255) NOT NULL,
    description TEXT,
    event_time TIMESTAMP NOT NULL,
    event_date DATE NOT NULL,
    location_id INTEGER REFERENCES Locations(location_id) NOT NULL,
    contact_phone VARCHAR(255),
    contact_email VARCHAR(255),
    visibility event_type NOT NULL,
    university_id INTEGER REFERENCES Universities(university_id),
    rso_id INTEGER REFERENCES RSOs(rso_id),
    super_admin_approved BOOLEAN DEFAULT FALSE
);
CREATE FUNCTION validate_event_creator() RETURNS TRIGGER AS $$ BEGIN IF NEW.user_id NOT IN (
    SELECT user_id
    FROM Users
    WHERE auth = 'admin' -- Use auth instead of user_type
) THEN RAISE EXCEPTION 'Only admins can create events!';
ELSIF NEW.event_type = 'public'
AND NEW.rso_id IS NULL THEN
UPDATE Events
SET super_admin_approved = FALSE
WHERE event_id = NEW.event_id;
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER validate_event_creation BEFORE
INSERT ON Events FOR EACH ROW EXECUTE PROCEDURE validate_event_creator();
CREATE TABLE Event_Feedback (
    feedback_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES Users(user_id) NOT NULL,
    event_id INTEGER REFERENCES Events(event_id) NOT NULL,
    content TEXT,
    -- Can store comment text or be null for ratings
    rating INTEGER CHECK (
        rating >= 1
        AND rating <= 5
    ),
    -- Rating value (null for comments)
    feedback_type VARCHAR(255) CHECK (feedback_type IN ('comment', 'rating')),
    -- Differentiates comment or rating
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- TODO: Debug this trigger
-- CREATE FUNCTION validate_rso_creation() RETURNS TRIGGER AS $$ BEGIN -- Declare variables
-- DECLARE user_count INTEGER;
-- DECLARE user_university INTEGER;
-- DECLARE is_admin_set BOOLEAN DEFAULT FALSE;
-- -- Execute SELECT statements and assign results to variables
-- SELECT COUNT(*) INTO user_count -- Use COUNT(*) instead of COUNT(user_id)
-- FROM User_RSO_Membership
-- WHERE rso_id = NEW.rso_id;
-- SELECT university_id INTO user_university
-- FROM Users
-- WHERE user_id = NEW.admin_id;
-- -- Rest of the function logic remains the same
-- IF user_count < 4 THEN RAISE EXCEPTION 'RSO creation requires at least 4 members!';
-- ELSIF user_university != (
--     SELECT university_id
--     FROM Users
--     WHERE user_id IN (
--             SELECT user_id
--             FROM User_RSO_Membership
--             WHERE rso_id = NEW.rso_id
--         )
-- ) THEN RAISE EXCEPTION 'All RSO members must be from the same university!';
-- ELSIF NEW.admin_id IS NULL THEN RAISE EXCEPTION 'An admin must be set for the RSO!';
-- END IF;
-- RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;
-- CREATE TRIGGER validate_rso_on_insert BEFORE
-- INSERT ON RSOs FOR EACH ROW EXECUTE PROCEDURE validate_rso_creation();