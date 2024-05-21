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

type ProductService struct {
	queries *db.Queries
	db      *pgxpool.Pool
	logger  *logger.Logger
}

func (p *ProductService) GetProductsCursor(ctx context.Context, name string, price float64, cursor string, size int) (res db_util.PageCursor, err error) {
	countProducts, err := p.queries.CountProducts(ctx, db.CountProductsParams{
		Name:  name,
		Price: price,
	})
	if err != nil {
		return
	}
	res.TotalData = int(countProducts)

	decodeString, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return
	}

	nextPrev, id, currPage, err := db_util.ParseCursor(string(decodeString))
	if err != nil {
		return
	}
	res.CurrentPage = currPage

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
				return
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
				return
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
			return
		}

		products = ConvertGetProductsCursorRowArrayToProductArray(p)
	}
	res.Data = products

	if len(products) == 0 {
		return res, errors.New("data not found")
	}

	res.TotalPage = db_util.CountPageSize(int(countProducts), size)

	if res.CurrentPage > res.TotalPage {
		return res, errors.New("page not found")
	}

	if res.TotalPage > res.CurrentPage {
		res.NextCursor = fmt.Sprintf("next,%d,%d", products[len(products)-1].ID, res.CurrentPage+1)
	}

	if res.CurrentPage > 1 {
		res.PrevCursor = fmt.Sprintf("prev,%d,%d", products[0].ID, res.CurrentPage-1)
	}

	res.NextCursor = base64.StdEncoding.EncodeToString([]byte(res.NextCursor))
	res.PrevCursor = base64.StdEncoding.EncodeToString([]byte(res.PrevCursor))

	return
}

func (p *ProductService) UpdateProductPrice(ctx context.Context, id int, name string, price float64) error {
	return p.queries.UpdateProduct(ctx, db.UpdateProductParams{
		Price: price,
		Name:  name,
		ID:    int32(id),
	})
}

func (p *ProductService) DeleteProduct(ctx context.Context, id int) error {
	return p.queries.DeleteProduct(ctx, int32(id))
}

func (p *ProductService) CreateProduct(ctx context.Context, name string, price float64) error {
	return p.queries.CreateProduct(ctx, db.CreateProductParams{
		Name:  name,
		Price: price,
	})
}

func (p *ProductService) GetProducts(ctx context.Context, name string, price float64, sortBy string, page, size int) (res db_util.Page, err error) {

	limit, offset, err := db_util.PageToLimitOffset(size, page)
	if err != nil {
		return
	}

	products, err := p.queries.GetProducts(ctx, db.GetProductsParams{
		Name:   fmt.Sprintf("%%%s%%", name),
		Price:  price,
		SortBy: sortBy,
		Offs:   int32(offset),
		Si:     int32(limit),
	})
	if err != nil {
		return
	}
	res.Data = products

	countProducts, err := p.queries.CountProducts(ctx, db.CountProductsParams{
		Name:  name,
		Price: price,
	})
	if err != nil {
		return
	}
	res.TotalData = int(countProducts)
	res.TotalPage = db_util.CountPageSize(int(countProducts), size)

	return
}

func ConvertGetProductsNextCursorRowArrayToProductArray(rows []db.GetProductsNextCursorRow) (res []model.Product) {
	for _, row := range rows {
		res = append(res, model.Product{
			ID:        row.ID,
			Name:      row.Name,
			Price:     row.Price,
			CreatedAt: row.CreatedAt,
		})
	}
	return
}

func ConvertGetProductsPrevCursorRowArrayToProductArray(rows []db.GetProductsPrevCursorRow) (res []model.Product) {
	for _, row := range rows {
		res = append(res, model.Product{
			ID:        row.ID,
			Name:      row.Name,
			Price:     row.Price,
			CreatedAt: row.CreatedAt,
		})
	}
	return
}

func ConvertGetProductsCursorRowArrayToProductArray(rows []db.GetProductsCursorRow) (res []model.Product) {
	for _, row := range rows {
		res = append(res, model.Product{
			ID:        row.ID,
			Name:      row.Name,
			Price:     row.Price,
			CreatedAt: row.CreatedAt,
		})
	}
	return
}

func NewProductService(
	queries *db.Queries,
	db *pgxpool.Pool,
	logger *logger.Logger,
) *ProductService {
	return &ProductService{
		queries: queries,
		db:      db,
		logger:  logger,
	}
}
