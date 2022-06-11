package dal

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Manuel9550/d20-workout/pkg/entities"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type DBManager struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

type ResourceNotFoundError struct {
	resourceType string
	resourceName string
}

type ResourceDuplicateError struct {
	resourceType string
	resourceName string
}

func (e ResourceNotFoundError) Error() string {
	return fmt.Sprintf("Could not find resource '%s' of type: '%s'", e.resourceName, e.resourceType)
}

func (e ResourceDuplicateError) Error() string {
	return fmt.Sprintf("The resource '%s' of type: '%s' already exists!", e.resourceName, e.resourceType)
}

func NewDBManager(connectionString string, logger *logrus.Logger) (*DBManager, error) {

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	// Test connection to database
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	newDBManager := DBManager{
		DB:     db,
		Logger: logger,
	}

	return &newDBManager, nil
}

func (dm *DBManager) GetUser(ctx context.Context, userName string) (*entities.User, error) {

	queryStatement := `SELECT UserName FROM D20WorkoutUser WHERE UserName = $1`

	foundUser := entities.User{}
	err := dm.DB.QueryRowContext(ctx, queryStatement, userName).Scan(&foundUser.Username)

	if err != nil {
		if err != sql.ErrNoRows {
			dm.Logger.WithFields(logrus.Fields{
				"QueryError": err,
				"Query":      queryStatement,
			}).Error()
			return nil, err
		}

		return nil, &ResourceNotFoundError{
			resourceName: userName,
			resourceType: "User",
		}
	}

	return &foundUser, nil
}

func (dm *DBManager) CreateUser(ctx context.Context, userName string) (*entities.User, error) {
	// Does user already exist?
	_, err := dm.GetUser(ctx, userName)
	if err != nil {
		return nil, err
	}

	insertionStatement := `INSERT INTO D20WorkoutUser(UserName) VALUES($1) RETURNING UserName`
	createdUser := entities.User{}
	err = dm.DB.QueryRowContext(ctx, insertionStatement, userName).Scan(&createdUser.Username)

	if err != nil {
		dm.Logger.WithFields(logrus.Fields{
			"QueryError": err,
			"Query":      insertionStatement,
		}).Error()
		return nil, err
	}

	return &createdUser, nil
}

func (dm *DBManager) CheckExerciseNumber(ctx context.Context, exerciseNumber int) error {
	queryStatement := `SELECT RollNumber FROM Exercise WHERE RollNumber = $1`

	foundExercise := entities.Exercise{}
	err := dm.DB.QueryRowContext(ctx, queryStatement, exerciseNumber).Scan(&foundExercise.RollNumber)

	if err != nil {
		if err != sql.ErrNoRows {
			dm.Logger.WithFields(logrus.Fields{
				"QueryError": err,
				"Query":      queryStatement,
			}).Error()
			return err
		}

		return &ResourceNotFoundError{
			resourceName: strconv.Itoa(exerciseNumber),
			resourceType: "RollNumber",
		}
	}

	return nil
}

func (dm *DBManager) ValidatePoint(ctx context.Context, point *entities.Point) error {
	// Does user exists?
	_, err := dm.GetUser(ctx, point.Username)
	if err != nil {
		return err
	}

	err = dm.CheckExerciseNumber(ctx, point.ExerciseNumber)
	if err != nil {
		return err
	}

	return nil
}

func (dm *DBManager) AddUserPoint(ctx context.Context, exercisePoint *entities.Point) error {
	err := dm.ValidatePoint(ctx, exercisePoint)
	if err != nil {
		dm.Logger.WithFields(logrus.Fields{
			"PointError": "Error When validating point",
			"Error":      err,
		}).Error()
		return err
	}

	insertionStatement := `INSERT INTO FinishedPoint(UserName,RollNumber,AmountDone,Timestamp) VALUES($1,$2,$3,$4)`
	_, err = dm.DB.ExecContext(ctx, insertionStatement, exercisePoint.Username, exercisePoint.ExerciseNumber, exercisePoint.Amount, exercisePoint.Timestamp)

	if err != nil {
		dm.Logger.WithFields(logrus.Fields{
			"QueryError": err,
			"Query":      insertionStatement,
		}).Error()
		return err
	}

	return nil
}

/*

func (dm *DBManager) SyncUser(ctx context.Context, exercisePoints entities.ExercisePoints) (*entities.User, error) {

	// Check that the user exists
	_, err := dm.GetUser(ctx, exercisePoints.Username)
	if err != nil {
		return nil, err
	}

	insertionStatement := `INSERT INTO D20WorkoutUser(UserName) VALUES($1) RETURNING UserName`
	createdUser := entities.User{}
	err = dm.DB.QueryRowContext(ctx, insertionStatement, userName).Scan(&createdUser.Username)

	if err != nil {
		dm.Logger.WithFields(logrus.Fields{
			"QueryError": err,
			"Query":      insertionStatement,
		}).Error()
		return nil, err
	}

	return &createdUser, nil
}
*/
