package postgres

import (
	"elkonterm/internal/repository/dto"
	"elkonterm/pkg/logger"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	categoriesTable = "categories"
)

type categoryRepo struct {
	db     *sqlx.DB
	logger *logger.Logger
}

func NewCategoryRepo(db *sqlx.DB, logger *logger.Logger) *categoryRepo {
	return &categoryRepo{db, logger}
}

func (r *categoryRepo) Create(input dto.CreateCategory) (uuid.UUID, error) {
	id := uuid.New()

	name, err := json.Marshal(input.Name)
	if err != nil {
		return uuid.Nil, err
	}

	query := sq.Insert(categoriesTable).
		Columns("id, parent_id, name, type, photo, is_top, application_areas_img, status").
		Values(id, input.ParentID, name, input.Type, input.Photo, input.IsTop, input.ApplicationAreasImg, input.Status).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		r.logger.Error(err)
		return uuid.Nil, err
	}

	stmt, err := r.db.Prepare(sqlQuery)
	if err != nil {
		r.logger.Error(err)
		return uuid.Nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		r.logger.Error(err)
		return uuid.Nil, err
	}

	return id, nil
}

func (r *categoryRepo) GetList(options dto.FilterOptions) ([]dto.Category, int, error) {
	categories := []dto.Category{}

	query := sq.
		Select("id, parent_id, name, type, photo, is_top, application_areas_img, status").
		From(categoriesTable).
		Where(sq.Eq{"deleted_at": nil})

	countQuery := sq.Select("count(id)").From(categoriesTable).Where(sq.Eq{"deleted_at": nil})

	filters := options.Filters

	if searchKey, ok := filters["search-key"]; ok {
		searchValue := "%" + searchKey.(string) + "%"
		query = query.Where(sq.Expr("name::varchar ILIKE ?", searchValue))
		countQuery = countQuery.Where(sq.Expr("name::varchar ILIKE ?", searchValue))
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

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		r.logger.Error(err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			category dto.Category
			name     []byte
		)

		if err = rows.Scan(
			&category.ID,
			&category.ParentID,
			&name,
			&category.Type,
			&category.Photo,
			&category.IsTop,
			&category.ApplicationAreasImg,
			&category.Status,
		); err != nil {
			r.logger.Error(err)
			return nil, 0, err
		}

		err = json.Unmarshal(name, &category.Name)
		if err != nil {
			r.logger.Error(err)
			return nil, 0, err
		}

		categories = append(categories, category)
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

	return categories, totalCount, nil
}
