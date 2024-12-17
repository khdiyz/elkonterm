package postgres

import (
	"database/sql"
	"elkonterm/internal/repository/dto"
	"elkonterm/pkg/logger"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
)

type UserRepo struct {
	db     *sqlx.DB
	logger *logger.Logger
}

func NewUserRepo(db *sqlx.DB, logger *logger.Logger) *UserRepo {
	return &UserRepo{db, logger}
}

func (r *UserRepo) Create(input dto.CreateUser) (uuid.UUID, error) {
	id := uuid.New()

	query, args, err := sq.Insert(usersTable).
		Columns("full_name, phone_number, role_id, email, password, company, status").
		Values(input.FullName, input.PhoneNumber, input.RoleID, input.Email, input.Password, input.Company, input.Status).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		r.logger.Error(err)
		return uuid.Nil, err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		r.logger.Error(err)
		return uuid.Nil, err
	}

	return id, nil
}

func (r *UserRepo) GetList(options dto.FilterOptions) ([]dto.User, int, error) {
	users := []dto.User{}

	query := sq.
		Select("id, full_name, phone_number, role_id, email, password, company, status, created_at").
		From(usersTable).
		Where(sq.Eq{"deleted_at": nil})

	countQuery := sq.Select("COUNT(id)").From(usersTable).Where(sq.Eq{"deleted_at": nil})

	filters := options.Filters

	if searchKey, ok := filters["search-key"]; ok {
		query = query.Where(sq.ILike{"(full_name || email || phone_number)": "%" + searchKey.(string) + "%"})
		countQuery = countQuery.Where(sq.ILike{"(full_name || email || phone_number)": "%" + searchKey.(string) + "%"})
	}

	if status, ok := filters["status"]; ok {
		query = query.Where(sq.Eq{"status": status})
		countQuery = countQuery.Where(sq.Eq{"status": status})
	}

	if options.SortBy != "" {
		order := "ASC"
		if options.Order == "desc" {
			order = "DESC"
		}
		query = query.OrderBy(fmt.Sprintf("%s %s", options.SortBy, order))
	} else {
		query = query.OrderBy("created_at DESC") // Default sorting
	}

	// Pagination
	if options.Limit > 0 {
		offset := (options.Page - 1) * options.Limit
		query = query.Limit(uint64(options.Limit)).Offset(uint64(offset))
	}

	// Build and execute the main query
	sqlQuery, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		r.logger.Error(err)
		return nil, 0, err
	}

	err = r.db.Select(&users, sqlQuery, args...)
	if err != nil {
		r.logger.Error(err)
		return nil, 0, err
	}

	// Build and execute the count query
	countSql, countArgs, err := countQuery.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		r.logger.Error(err)
		return nil, 0, err
	}

	var totalCount int
	err = r.db.Get(&totalCount, countSql, countArgs...)
	if err != nil {
		r.logger.Error(err)
		return nil, 0, err
	}

	return users, totalCount, nil
}

func (r *UserRepo) GetByEmail(email string) (dto.User, error) {
	var user dto.User

	query, args, err := sq.
		Select("id, full_name, phone_number, role_id, email, password, company, status, created_at").
		From(usersTable).
		Where(sq.And{
			sq.Eq{"email": email},
			sq.Eq{"deleted_at": nil},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		r.logger.Error(err)
		return user, err
	}

	err = r.db.Get(&user, query, args...)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			r.logger.Error(err)
		}
		return user, err
	}

	return user, nil
}
