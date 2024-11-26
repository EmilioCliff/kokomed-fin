package mysql

import (
	"context"
	"database/sql"
	"log"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Store struct {
	db     *sql.DB
	config pkg.Config
}

func NewStore(config pkg.Config) *Store {
	return &Store{
		config: config,
	}
}

type MySQLRepo struct {
	Branches  repository.BranchRepository
	Loans     repository.LoansRepository
	Products  repository.ProductRepository
	Users     repository.UserRepository
	Clients   repository.ClientRepository
	NonPosted repository.NonPostedRepository
}

func NewMySQLRepo(db *Store) *MySQLRepo {
	return &MySQLRepo{
		Branches:  NewBranchRepository(db),
		Loans:     NewLoanRepository(db),
		Products:  NewProductRepository(db),
		Users:     NewUserRepository(db),
		Clients:   NewClientRepository(db),
		NonPosted: NewNonPostedRepository(db),
	}
}

// open db.
func (s *Store) OpenDB() error {
	if s.config.DB_DSN == "" {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "dsn is empty")
	}

	var err error

	s.db, err = sql.Open("mysql", s.config.DB_DSN)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to open database: %v", err)
	}

	err = s.db.Ping()
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to ping database: %v", err)
	}

	driver, err := mysql.WithInstance(s.db, &mysql.Config{})
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create driver: %v", err)
	}

	return s.runMigration(driver)
}

// close db.
func (s *Store) CloseDB() error {
	log.Println("Shutting down database...")

	if s.db != nil {
		return s.db.Close()
	}

	return nil
}

// run migration.
func (s *Store) runMigration(driver database.Driver) error {
	if s.config.MIGRATION_PATH == "" {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "migrations directory is empty")
	}

	migration, err := migrate.NewWithDatabaseInstance(s.config.MIGRATION_PATH, s.config.DB_DSN, driver)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "Failed to load migration: %s", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "Failed to run migrate up: %s", err)
	}

	return nil
}

// executes transaction.
func (s *Store) ExecTx(ctx context.Context, fn func(q *generated.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "couldn't being transaction: %v", err)
	}

	q := generated.New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "tx err: %v, rb err: %v", err, rbErr)
		}

		return pkg.Errorf(pkg.INTERNAL_ERROR, "tx err: %v", err)
	}

	return tx.Commit()
}
