package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/danangkonang/cake-store-restful/config"
	"github.com/danangkonang/cake-store-restful/model"
)

type ServiceCake interface {
	FindCakes() ([]*model.ProductResponse, error)
	FindCakeById(id int) (*model.ProductResponse, error)
	SaveCake(ck *model.ProductPostRequest) error
	UpdateCake(ck *model.ProductUpdateRequest) error
	DeleteCake(ck *model.ProductDeleteRequest) error
}

func NewServiceCake(Con *config.DB) ServiceCake {
	return &connection{
		Mysql: Con.Mysql,
	}
}

func (r *connection) FindCakes() ([]*model.ProductResponse, error) {
	fac := make([]*model.ProductResponse, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "SELECT id, title, description, rating, image, created_at, updated_at FROM products"
	row, err := r.Mysql.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		c := new(model.ProductResponse)
		var update sql.NullTime
		var created time.Time
		err := row.Scan(&c.Id, &c.Title, &c.Description, &c.Rating, &c.Image, &created, &update)
		if err != nil {
			return nil, err
		}
		c.CreatedAt = created.Format("2006-01-02 15:04:05")
		if update.Valid {
			c.UpdatedAt = update.Time.Format("2006-01-02 15:04:05")
		}
		fac = append(fac, c)
	}
	return fac, nil
}

func (r *connection) SaveCake(c *model.ProductPostRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO products (title, description, rating, image, created_at) VALUES (?, ?, ?, ?, ?)"
	_, err := r.Mysql.ExecContext(ctx, query, c.Title, c.Description, c.Rating, c.Image, c.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *connection) UpdateCake(c *model.ProductUpdateRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "UPDATE products SET title=?, description=?, rating=?, image=?, updated_at=? WHERE id=?"
	_, err := r.Mysql.ExecContext(ctx, query, c.Title, c.Description, c.Rating, c.Image, time.Now(), c.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *connection) DeleteCake(c *model.ProductDeleteRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "DELETE FROM products WHERE id=?"
	_, err := r.Mysql.ExecContext(ctx, query, c.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *connection) FindCakeById(id int) (*model.ProductResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "SELECT id, title, description, rating, image, created_at, updated_at FROM products WHERE id=?"
	row := r.Mysql.QueryRowContext(ctx, query, id)
	c := new(model.ProductResponse)
	var update sql.NullTime
	var created time.Time
	err := row.Scan(&c.Id, &c.Title, &c.Description, &c.Rating, &c.Image, &created, &update)
	if err != nil {
		return nil, err
	}
	c.CreatedAt = created.Format("2006-01-02 15:04:05")
	if update.Valid {
		c.UpdatedAt = update.Time.Format("2006-01-02 15:04:05")
	}
	return c, nil
}
