CREATE TABLE Location (
    Lname VARCHAR(255) PRIMARY KEY,
    Address TEXT,
    Longitude REAL,
    Latitude REAL
);
CREATE TABLE Users (
    UID INTEGER PRIMARY KEY,
    Username VARHAR(255),
    firstName VARCHAR(255),
    lastName VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    Location TEXT REFERENCES Location(Lname)
);
CREATE TABLE RSO_Events (Event_ID INTEGER PRIMARY KEY, RSOS TEXT);
CREATE TABLE Events (
    Event_ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Time DATETIME,
    Location TEXT REFERENCES Location(Lname),
    Event_name TEXT,
    Description TEXT,
    ISA TEXT,
    CONSTRAINT Event_Unique UNIQUE (Time, Location)
);
CREATE TABLE Comments (
    Event_ID INTEGER,
    timestamp DATETIME,
    rating INTEGER,
    comment TEXT,
    FOREIGN KEY (Event_ID) REFERENCES Events(Event_ID)
);
CREATE TABLE Private_Events (
    Event_ID INTEGER PRIMARY KEY,
    FOREIGN KEY (Event_ID) REFERENCES Events(Event_ID),
    CONSTRAINT FK_Private_Event_Unique UNIQUE (Event_ID)
);
CREATE TABLE Public_Events (
    Event_ID INTEGER PRIMARY KEY,
    FOREIGN KEY (Event_ID) REFERENCES Events(Event_ID),
    CONSTRAINT FK_Public_Event_Unique UNIQUE (Event_ID)
);
CREATE TABLE Admins (
    UID INTEGER,
    FOREIGN KEY (UID) REFERENCES Users(UID)
);
CREATE TABLE SuperAdmins (
    UID INTEGER,
    FOREIGN KEY (UID) REFERENCES Users(UID)
);
-- Ternary Relationships (Many-to-Many)
CREATE TABLE Creates_Public_Events (
    Event_ID INTEGER REFERENCES Public_Events(Event_ID),
    Admin_UID INTEGER REFERENCES Admins(UID),
    SuperAdmin_UID INTEGER REFERENCES SuperAdmins(UID),
    PRIMARY KEY (Event_ID, Admin_UID, SuperAdmin_UID)
);
CREATE TABLE Creates_Private_Events (
    Event_ID INTEGER REFERENCES Private_Events(Event_ID),
    Admin_UID INTEGER REFERENCES Admins(UID),
    SuperAdmin_UID INTEGER REFERENCES SuperAdmins(UID),
    PRIMARY KEY (Event_ID, Admin_UID, SuperAdmin_UID)
);
CREATE TABLE Universities (
    University_ID INTEGER PRIMARY KEY,
    Name TEXT
);
CREATE TABLE Profiles (
    UID INTEGER PRIMARY KEY,
    University_ID INTEGER,
    FOREIGN KEY (UID) REFERENCES Users(UID),
    FOREIGN KEY (University_ID) REFERENCES Universities(University_ID)
);
-- Foreign Key Constraints (ON DELETE CASCADE)
ALTER TABLE Events
ADD CONSTRAINT FK_Events_Creates_Users FOREIGN KEY (Creates) REFERENCES Users(UID) ON DELETE CASCADE;
ALTER TABLE RSO_Events
ADD CONSTRAINT FK_RSO_Events_Owns_Users FOREIGN KEY (Owns) REFERENCES Users(UID) ON DELETE CASCADE;
ALTER TABLE Comments
ADD CONSTRAINT FK_Comments_At_Users FOREIGN KEY (At) REFERENCES Users(UID) ON DELETE CASCADE;
ALTER TABLE Public_Events
ADD CONSTRAINT FK_Public_Events_Creates_Public FOREIGN KEY (Event_ID) REFERENCES Creates_Public_Events(Event_ID) ON DELETE CASCADE;
ALTER TABLE Private_Events
ADD CONSTRAINT FK_Private_Events_Creates_Private FOREIGN KEY (Event_ID) REFERENCES Creates_Private_Events(Event_ID) ON DELETE CASCADE;
-- Implement Constraints as CHECK constraints (Assertions in Postgres)
ALTER TABLE Events
ADD CHECK (
        ISA IN ('RSO_Events', 'Private Events', 'Public Events')
    );
CREATE ASSERTION is_event_subclass CHECK (
    NOT EXISTS (
        SELECT 1
        FROM RSO_Events
        WHERE Event_ID NOT IN (
                SELECT Event_ID
                FROM Events
            )
    )
);
-- RSO_Events ⊆ Events
CREATE ASSERTION disjoint_private_public CHECK (
    NOT EXISTS (
        SELECT 1
        FROM Private_Events
        WHERE Event_ID IN (
                SELECT Event_ID
                FROM Public_Events
            )
    )
);
-- RSO_Events ∩ Private_Events = Ø
CREATE ASSERTION covering_events CHECK (
    NOT EXISTS (
        SELECT 1
        FROM Events
        WHERE Event_ID NOT IN (
                SELECT Event_ID
                FROM RSO_Events
                UNION ALL
                SELECT Event_ID
                FROM Private_Events
                UNION ALL
                SELECT Event_ID
                FROM Public_Events
            )
    )
);
-- RSO_Events ∪ Private_Events ∪ Public_Events = Events