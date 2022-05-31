
DROP TABLE IF EXISTS FinishedPoint;
DROP TABLE IF EXISTS D20WorkoutUser;
DROP TABLE IF EXISTS Exercise;

-- Create the tables
CREATE TABLE IF NOT EXISTS D20WorkoutUser (
    UserName VARCHAR NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS Exercise (
    RollNumber INTEGER NOT NULL PRIMARY KEY,
    ExerciseName VARCHAR NOT NULL,
    StartingAmount INTEGER NOT NULL,
    IncrementAmount INTEGER NOT NULL,
    Units VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS FinishedPoint (
    ID SERIAL,
    UserName VARCHAR REFERENCES D20WorkoutUser (UserName),
    RollNumber INTEGER REFERENCES Exercise(Rollnumber),
    AmountDone INTEGER NOT NULL,
    Timestamp TIMESTAMP NOT NULL
);

-- Set up the initial data (TODO: Find a way to customize this? Put it into another table, 'contest', that has a start and end date?)
INSERT INTO Exercise(RollNumber, ExerciseName, StartingAmount, IncrementAmount, Units)
VALUES (1, 'Burpees', 100, 0,'Sets'),
(2, 'Horsestance', 45, 5,'Sets'),
(3, 'Hang', 30, 2,'Seconds'),
(4, 'Hollow Body', 45, 5,'Seconds'),
(5, 'Flutter Kicks', 100, 10,'Sets'),
(6, 'Raised Leg Lunges', 10, 2,'Sets'),
(7, 'Dips on Chair', 10, 2,'Sets'),
(8, 'L-Sit', 30, 2,'Seconds'),
(9, 'Bear Crawls', 20, 2,'Sets'),
(10, 'Wildcard', 0, 0,'Users Choice'),
(11, 'Glute Bridges', 60, 5,'Sets'),
(12, 'Calf Raises', 30, 2,'Sets'),
(13, 'Wall Sits', 30, 2,'Seconds'),
(14, 'Pushups', 20, 2,'Sets'),
(15, 'Plank', 45, 5,'Seconds'),
(16, 'Side Plank', 30, 2,'Seconds'),
(17, 'Superman', 45, 5,'Seconds'),
(18, 'Jumping Jacks', 100, 10,'Sets'),
(19, 'Pullups', 2, 1,'Sets'),
(20, 'Freebie', 0, 0,'free point');

