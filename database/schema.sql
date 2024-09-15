CREATE TABLE CLASSES (
    ID SERIAL PRIMARY KEY,
    NAME VARCHAR(255) NOT NULL
);

CREATE TABLE OBJECTS (
    ID SERIAL PRIMARY KEY,
    NAME VARCHAR(255) NOT NULL,
    VALUE DECIMAL NOT NULL,
    DATETIME TIMESTAMP NOT NULL,
    CLASS_ID INTEGER REFERENCES CLASSES(ID) ON DELETE CASCADE
);

CREATE TABLE ITEMS (
    ID SERIAL PRIMARY KEY,
    NAME VARCHAR(255) NOT NULL,
    VALUE DECIMAL NOT NULL,
    DATETIME TIMESTAMP NOT NULL
);

CREATE TABLE OBJECT_ITEM_LINK (
    OBJECT_ID INTEGER REFERENCES OBJECTS(ID) ON DELETE CASCADE,
    ITEM_ID INTEGER REFERENCES ITEMS(ID) ON DELETE CASCADE,
    PRIMARY KEY (OBJECT_ID, ITEM_ID)
);
