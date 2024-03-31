-- Create Authority type
CREATE TYPE Authority AS ENUM ('normal', 'admin', 'superadmin');
-- Create Users table
CREATE TABLE Users (
    UID SERIAL PRIMARY KEY,
    Name TEXT NOT NULL,
    Email TEXT NOT NULL UNIQUE,
    Login TEXT NOT NULL UNIQUE,
    Password TEXT NOT NULL,
    DateCreated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Desc TEXT NOT NULL,
    Authority Authority NOT NULL DEFAULT 'normal'
);
-- Create Comments table
CREATE TABLE Comments (
    Text TEXT NOT NULL,
    rating INT,
    timestamp TIMESTAMP NOT NULL,
    UID INT REFERENCES Users(UID)
);
-- Create Location table
CREATE TABLE Location (
    Lname TEXT PRIMARY KEY,
    Address TEXT,
    Longitude FLOAT,
    Latitude FLOAT
);
-- Create Events table
CREATE TABLE Events (
    EventID SERIAL PRIMARY KEY,
    Time TIME NOT NULL,
    Lname TEXT REFERENCES Location(Lname),
    Event_name TEXT,
    Description TEXT,
    UID INT REFERENCES Users(UID),
    CommentID INT REFERENCES Comments(CommentID),
    UNIQUE (Time, Lname)
);
-- Create RSOs table
CREATE TABLE RSOs (RSO_ID SERIAL PRIMARY KEY);
-- Create Admins table
CREATE TABLE Admins (
    Admin_ID SERIAL PRIMARY KEY,
    UID INT REFERENCES Users(UID)
);
-- Create SuperAdmins table
CREATE TABLE SuperAdmins (
    SuperAdmin_ID SERIAL PRIMARY KEY,
    UID INT REFERENCES Users(UID)
);
-- Create RSO_Events table and establish inheritance relationship with Events table
CREATE TABLE RSO_Events () INHERITS (Events);
-- Add foreign key constraint for ownership between RSO_Events and RSOs tables 
ALTER TABLE RSO_Events
ADD CONSTRAINT fk_rso_events_rso_id FOREIGN KEY (RSO_ID) REFERENCES RSOs(RSO_ID);
-- Create Private_Events table and establish inheritance relationship with Events table 
CREATE TABLE Private_Events (
    Admin_ID INT NOT NULL,
    SuperAdmin_ID INT NOT NULL,
    FOREIGN KEY (Admin_ID) REFERENCES Admins(Admin_ID),
    FOREIGN KEY (SuperAdmin_ID) REFERENCES SuperAdmins(SuperAdmin_ID)
) INHERITS (Events);
-- Create Public_Events table and establish inheritance relationship with Events Table 
CREATE TABLE Public_Events (
    Admin_ID INT NOT NULL,
    SuperAdmin_ID INT NOT NULL,
    FOREIGN KEY (Admin_ID) REFERENCES Admins(Admin_ID),
    FOREIGN KEY (SuperAdmin_ID) REFERENCES SuperAdmins(SuperAdmin_ID)
) INHERITS (Events);
-- Add check constraints for ISA, disjointness and covering
ALTER TABLE RSO_Events
ADD CONSTRAINT check_isa CHECK (RSO_ID IS NOT NULL);
ALTER TABLE Private_Events
ADD CONSTRAINT check_disjointness CHECK (RSO_ID IS NULL);
ALTER TABLE Events
ADD CONSTRAINT check_covering CHECK (
        RSO_ID IS NOT NULL
        OR Admin_ID IS NOT NULL
        OR SuperAdmin_ID IS NOT NULL
    );