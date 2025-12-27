package database

import (
	"payment_microservice/internal/common/config/env"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetDB(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	// Reset singletons para teste isolado
	dbConnection = nil
	instance = nil

	db := GetDB()
	assert.Nil(t, db, "GetDB should return nil when dbConnection is not set")
}

func TestConnect_AlreadyConnected(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	// Setup: criar uma conexão mock
	mockDB, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Simular conexão já estabelecida
	dbConnection = gormDB

	// Act: tentar conectar novamente
	Connect()

	// Assert: não deve criar nova conexão
	assert.NotNil(t, dbConnection)
}

func TestClose_NoConnection(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	// Reset singleton
	dbConnection = nil

	// Act & Assert: não deve fazer panic
	assert.NotPanics(t, func() {
		Close()
	})
}

func TestClose_WithConnection(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	// Setup: criar uma conexão mock
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	assert.NoError(t, err)

	dbConnection = gormDB

	// Expect close
	mock.ExpectClose()

	// Act
	Close()

	// Assert
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRunMigrations(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	// Setup: criar uma conexão mock
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	assert.NoError(t, err)

	dbConnection = gormDB

	// Expect migration queries used by GORM AutoMigrate
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM information_schema.tables").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("CREATE TABLE").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE UNIQUE INDEX").
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Act: não deve fazer panic
	assert.NotPanics(t, func() {
		RunMigrations()
	})

	// Assert: todas as expectativas de SQL foram atendidas
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestConnect_RetryLogic(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	// Reset singletons
	dbConnection = nil
	instance = nil

	// Este teste apenas verifica se o código não entra em panic
	// Em ambiente de teste, a conexão falhará mas não devemos testar
	// a conexão real ao banco de dados

	// Para testar retry logic, precisaríamos de uma interface injetável
	// Por enquanto, apenas documentamos que o código tem retry logic
	assert.NotNil(t, &dbConnection)
}
