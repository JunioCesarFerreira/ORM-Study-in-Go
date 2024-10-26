-- Table for Projects
CREATE TABLE PROJECTS (
    ID          SERIAL PRIMARY KEY,
    NAME        VARCHAR(255) NOT NULL,
    MANAGER     VARCHAR(255) NOT NULL,    
    START_DATE  DATE NOT NULL,
    END_DATE    DATE,
    BUDGET      DECIMAL(12, 2),   
    DESCRIPTION TEXT               
);

-- Table for Tasks
CREATE TABLE TASKS (
    ID             SERIAL PRIMARY KEY,
    NAME           VARCHAR(255) NOT NULL,
    RESPONSIBLE    VARCHAR(255), 
    DEADLINE       DATE NOT NULL,           -- Deadline for task completion
    STATUS         VARCHAR(50) NOT NULL,    -- Task status (e.g., "in progress", "completed")
    PRIORITY       VARCHAR(20),             -- Task priority (e.g., "high", "medium", "low")
    ESTIMATED_TIME INTERVAL,          
    PROJECT_ID     INTEGER REFERENCES PROJECTS(ID) ON DELETE CASCADE,  
    DESCRIPTION    TEXT                 
);

-- Table for Resources
CREATE TABLE RESOURCES (
    ID               SERIAL PRIMARY KEY,
    TYPE             VARCHAR(255) NOT NULL,
    NAME             VARCHAR(255) NOT NULL,
    DAILY_COST       DECIMAL(10, 2),        
    STATUS           VARCHAR(50) NOT NULL,
    SUPPLIER         VARCHAR(255),
    QUANTITY         INTEGER, 
    ACQUISITION_DATE DATE
);

-- Link table between Tasks and Resources
CREATE TABLE TASK_RESOURCE (
    TASK_ID       INTEGER REFERENCES TASKS(ID) ON DELETE CASCADE,
    RESOURCE_ID   INTEGER REFERENCES RESOURCES(ID) ON DELETE CASCADE,
    QUANTITY_USED INTEGER,
    PRIMARY KEY (TASK_ID, RESOURCE_ID)
);
