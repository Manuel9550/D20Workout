package dal

import (
	"context"
	"database/sql"
	"fmt"

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

func (e ResourceNotFoundError) Error() string {
	return fmt.Sprintf("Could not find resource '%s' of type: '%s'", e.resourceName, e.resourceType)
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

	insertionStatement := `"INSERT INTO D20WorkoutUser(UserName) VALUES($1) RETURNING UserName`
	createdUser := entities.User{}
	err := dm.DB.QueryRowContext(ctx, insertionStatement, userName).Scan(createdUser.Username)

	if err != nil {
		dm.Logger.WithFields(logrus.Fields{
			"QueryError": err,
			"Query":      insertionStatement,
		}).Error()
		return nil, err
	}

	return &createdUser, nil
}
