package user

import (
	"context"
	"database/sql"
	"fmt"
	"pii-encrypt-example/entity"
	"pii-encrypt-example/pkg/exception"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	BeginTx(ctx context.Context) (tx *sql.Tx, err error)
	RollbackTx(ctx context.Context, tx *sql.Tx) (err error)
	CommitTx(ctx context.Context, tx *sql.Tx) (err error)
	SaveUser(ctx context.Context, user entity.User, tx *sql.Tx) (id int64, err error)
	// UpdateById(ctx context.Context, id int64, user UserRequest, tx *sql.Tx) (err error)
	FindManyUser(ctx context.Context, filter UserFilter) (bunchOfUsers []entity.User, err error)
	// FindOneUserByUUID(ctx context.Context, uuid string) (user entity.User, err error)
}

type sqlCommand interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type userRepository struct {
	logger      *logrus.Logger
	dbReadOnly  *sql.DB
	dbReadWrite *sql.DB
	tableName   string
}

// NewUserRepository is a constructor
func NewUserRepository(logger *logrus.Logger, dbReadOnly *sql.DB, dbReadWrite *sql.DB, tableName string) UserRepository {
	return &userRepository{
		logger:      logger,
		dbReadOnly:  dbReadOnly,
		dbReadWrite: dbReadWrite,
		tableName:   tableName,
	}
}

// BeginTx returns sql trx for global scope.
func (r *userRepository) BeginTx(ctx context.Context) (tx *sql.Tx, err error) {
	return r.dbReadWrite.BeginTx(ctx, nil)
}

// CommitTx will commit the transaction that has began.
func (r *userRepository) CommitTx(ctx context.Context, tx *sql.Tx) (err error) {
	return tx.Commit()
}

// RollbackTx will rollback the transaction to achieve the consistency.
func (r *userRepository) RollbackTx(ctx context.Context, tx *sql.Tx) (err error) {
	return tx.Rollback()
}

// Save will collect the order
func (r *userRepository) SaveUser(ctx context.Context, user entity.User, tx *sql.Tx) (id int64, err error) {
	var cmd sqlCommand = r.dbReadWrite
	if tx != nil {
		cmd = tx
	}

	command := fmt.Sprintf(`INSERT INTO %s SET uuid = ?, __encrypted__data_nama_crypt = ?, __encrypted__data_nama_hash = ?, __encrypted__data_email_crypt = ?, created_at = ?`, r.tableName)
	_, err = r.exec(ctx, cmd, command, user.UUID, user.NameCrypt, user.NameHash, user.EmailCrypt, user.CreatedAt)
	if err != nil {
		err = wrapError(err)
		return
	}

	// id, err = res.LastInsertId()
	// if err != nil {
	// 	err = wrapError(err)
	// 	return
	// }

	return
}

func (r *userRepository) FindManyUser(ctx context.Context, filter UserFilter) (bunchOfUsers []entity.User, err error) {
	var cmd sqlCommand = r.dbReadOnly
	var params []interface{}

	q := fmt.Sprintf(`SELECT u.uuid, u.__encrypted__data_nama_crypt, u.__encrypted__data_nama_hash, u.__encrypted__data_email_crypt, u.created_at FROM %s u`, r.tableName)

	if filter.Name != "" {
		q += fmt.Sprintf(` WHERE %s = ?`, "u.__encrypted__data_nama_hash")
		params = append(params, filter.NameHashed)
	}

	bunchOfUsers, err = r.query(ctx, cmd, q, params...)
	if err != nil {
		err = wrapError(err)
		return
	}
	return
}

func (r *userRepository) query(ctx context.Context, cmd sqlCommand, query string, args ...interface{}) (bunchOfUsers []entity.User, err error) {
	var rows *sql.Rows
	if rows, err = cmd.QueryContext(ctx, query, args...); err != nil {
		r.logger.WithContext(ctx).Error(query, err)
		return
	}

	defer func() {
		if err := rows.Close(); err != nil {
			r.logger.WithContext(ctx).Error(query, err)
		}
	}()

	for rows.Next() {
		var user entity.User

		err = rows.Scan(&user.UUID, &user.NameCrypt, &user.NameHash, &user.EmailCrypt, &user.CreatedAt)

		if err != nil {
			r.logger.WithContext(ctx).Error(query, err)
			return
		}

		bunchOfUsers = append(bunchOfUsers, user)
	}

	return
}

func (r *userRepository) exec(ctx context.Context, cmd sqlCommand, command string, args ...interface{}) (result sql.Result, err error) {
	var stmt *sql.Stmt
	if stmt, err = cmd.PrepareContext(ctx, command); err != nil {
		r.logger.WithContext(ctx).Error(command, err)
		return
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			r.logger.WithContext(ctx).Error(command, err)
		}
	}()

	if result, err = stmt.ExecContext(ctx, args...); err != nil {
		r.logger.WithContext(ctx).Error(command, err)
	}

	return
}

func wrapError(e error) (err error) {
	if e == sql.ErrNoRows {
		return exception.ErrNotFound
	}
	if driverErr, ok := e.(*mysql.MySQLError); ok {
		if driverErr.Number == 1062 {
			return exception.ErrConflict
		}
	}
	return exception.ErrInternalServer
}
