package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	db "github.com/ciazhar/go-zhar/examples/postgres/crud/internal/generated/repository"
	"github.com/ciazhar/go-zhar/examples/postgres/crud/internal/model"
	"github.com/ciazhar/go-zhar/pkg/db_util"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
)

type ProductService interface {
	CreateProduct(ctx context.Context, name string, price float64) error
	GetProducts(ctx context.Context, name string, price float64, sortBy string, page, size int) (db_util.Page, error)
	GetProductsCursor(ctx context.Context, name string, price float64, cursor string, size int) (db_util.PageCursor, error)
	UpdateProductPrice(ctx context.Context, id int, name string, price float64) error
	DeleteProduct(ctx context.Context, id int) error
}

type productService struct {
	queries *db.Queries
	db      *pgxpool.Pool
	logger  logger.Logger
}

func (p productService) GetProductsCursor(ctx context.Context, name string, price float64, cursor string, size int) (db_util.PageCursor, error) {
	countProducts, err := p.queries.CountProducts(ctx, db.CountProductsParams{
		Name:  name,
		Price: price,
	})
	if err != nil {
		return db_util.PageCursor{}, err
	}

	decodeString, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return db_util.PageCursor{}, err
	}

	nextPrev, id, currPage, err := db_util.ParseCursor(string(decodeString))
	if err != nil {
		return db_util.PageCursor{}, err
	}

	var products []model.Product
	switch nextPrev {
	case "next":
		fallthrough
	case "prev":
		idI, err := strconv.Atoi(id)
		if err != nil {
			return db_util.PageCursor{}, err
		}
		if nextPrev == "next" {
			p, err := p.queries.GetProductsNextCursor(ctx, db.GetProductsNextCursorParams{
				Name:   fmt.Sprintf("%%%s%%", name),
				Price:  price,
				Cursor: int32(idI),
				Si:     int32(size),
			})
			if err != nil {
				return db_util.PageCursor{}, err
			}

			products = ConvertGetProductsNextCursorRowArrayToProductArray(p)

		} else {
			p, err := p.queries.GetProductsPrevCursor(ctx, db.GetProductsPrevCursorParams{
				Cursor: int32(idI),
				Name:   fmt.Sprintf("%%%s%%", name),
				Price:  price,
				Si:     int32(size),
			})
			if err != nil {
				return db_util.PageCursor{}, err
			}

			products = ConvertGetProductsPrevCursorRowArrayToProductArray(p)
		}
	case "":
		fallthrough
	default:
		p, err := p.queries.GetProductsCursor(ctx, db.GetProductsCursorParams{
			Name:  fmt.Sprintf("%%%s%%", name),
			Price: price,
			Si:    int32(size),
		})

		if err != nil {
			return db_util.PageCursor{}, err
		}

		products = ConvertGetProductsCursorRowArrayToProductArray(p)
	}

	if len(products) == 0 {
		return db_util.PageCursor{}, errors.New("data not found")
	}

	totalPage := db_util.CountPageSize(int(countProducts), size)

	if currPage > totalPage {
		return db_util.PageCursor{}, errors.New("page not found")
	}

	nextCursor := ""
	if totalPage > currPage {
		nextCursor = fmt.Sprintf("next,%d,%d", products[len(products)-1].ID, currPage+1)
	}

	prevCursor := ""
	if currPage > 1 {
		prevCursor = fmt.Sprintf("prev,%d,%d", products[0].ID, currPage-1)
	}

	nextCursor = base64.StdEncoding.EncodeToString([]byte(nextCursor))
	prevCursor = base64.StdEncoding.EncodeToString([]byte(prevCursor))

	return db_util.PageCursor{
		Data:        products,
		TotalData:   int(countProducts),
		CurrentPage: currPage,
		TotalPage:   totalPage,
		NextCursor:  nextCursor,
		PrevCursor:  prevCursor,
	}, nil
}

func (p productService) UpdateProductPrice(ctx context.Context, id int, name string, price float64) error {
	return p.queries.UpdateProduct(ctx, db.UpdateProductParams{
		Price: price,
		Name:  name,
		ID:    int32(id),
	})
}

func (p productService) DeleteProduct(ctx context.Context, id int) error {
	return p.queries.DeleteProduct(ctx, int32(id))
}

func (p productService) CreateProduct(ctx context.Context, name string, price float64) error {
	return p.queries.CreateProduct(ctx, db.CreateProductParams{
		Name:  name,
		Price: price,
	})
}

func (p productService) GetProducts(ctx context.Context, name string, price float64, sortBy string, page, size int) (db_util.Page, error) {

	limit, offset, err := db_util.PageToLimitOffset(size, page)
	if err != nil {
		return db_util.Page{}, err
	}

	products, err := p.queries.GetProducts(ctx, db.GetProductsParams{
		Name:   fmt.Sprintf("%%%s%%", name),
		Price:  price,
		SortBy: sortBy,
		Offs:   int32(offset),
		Si:     int32(limit),
	})
	if err != nil {
		return db_util.Page{}, err
	}

	countProducts, err := p.queries.CountProducts(ctx, db.CountProductsParams{
		Name:  name,
		Price: price,
	})
	if err != nil {
		return db_util.Page{}, err
	}

	return db_util.Page{
		Data:      products,
		TotalData: int(countProducts),
		TotalPage: db_util.CountPageSize(int(countProducts), size),
	}, nil
}

func ConvertGetProductsNextCursorRowArrayToProductArray(rows []db.GetProductsNextCursorRow) []model.Product {
	var result []model.Product
	for _, row := range rows {
		result = append(result, model.Product{
			ID:        row.ID,
			Name:      row.Name,
			Price:     row.Price,
			CreatedAt: row.CreatedAt,
		})
	}
	return result
}

func ConvertGetProductsPrevCursorRowArrayToProductArray(rows []db.GetProductsPrevCursorRow) []model.Product {
	var result []model.Product
	for _, row := range rows {
		result = append(result, model.Product{
			ID:        row.ID,
			Name:      row.Name,
			Price:     row.Price,
			CreatedAt: row.CreatedAt,
		})
	}
	return result
}

func ConvertGetProductsCursorRowArrayToProductArray(rows []db.GetProductsCursorRow) []model.Product {
	var result []model.Product
	for _, row := range rows {
		result = append(result, model.Product{
			ID:        row.ID,
			Name:      row.Name,
			Price:     row.Price,
			CreatedAt: row.CreatedAt,
		})
	}
	return result
}

func NewProductService(
	queries *db.Queries,
	db *pgxpool.Pool,
	logger logger.Logger,
) ProductService {
	return &productService{
		queries: queries,
		db:      db,
		logger:  logger,
	}
}
