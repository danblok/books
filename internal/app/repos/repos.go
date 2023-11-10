package repos

import (
	"context"
	"errors"
	"log/slog"

	"github.com/danblok/books/internal/app/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrNoChange = errors.New("no rows were affected")

type BooksRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func New(db *sqlx.DB) *BooksRepo {
	return &BooksRepo{
		db:     db,
		logger: slog.Default(),
	}
}

func (r *BooksRepo) GetByID(ctx context.Context, id string) (*models.Book, error) {
	var book models.Book
	err := r.db.GetContext(ctx, &book, "SELECT * FROM Books WHERE id=$1", id)
	if err != nil {
		r.logger.Error("error in GetByID", err)
		return nil, err
	}
	return &book, nil
}

func (r *BooksRepo) GetAll(ctx context.Context) ([]*models.Book, error) {
	books := make([]*models.Book, 0)
	err := r.db.SelectContext(ctx, &books, "SELECT * FROM Books")
	if err != nil {
		r.logger.Error("error in GetAll", err)
		return nil, err
	}
	return books, nil
}

func (r *BooksRepo) Add(ctx context.Context, book models.Book) (*models.Book, error) {
	book.ID = uuid.NewString()
	res, err := r.db.ExecContext(ctx, "INSERT INTO Books (id, name, author, price) values ($1, $2, $3, $4)", book.ID, book.Name, book.Author, book.Price)
	if err != nil {
		r.logger.Error("error in Add", err)
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("error in Add", err)
		return nil, err
	}
	if rows == 0 {
		r.logger.Error("error in Add", ErrNoChange)
		return nil, ErrNoChange
	}
	return &book, nil
}

func (r *BooksRepo) Update(ctx context.Context, book models.Book) error {
	res, err := r.db.ExecContext(ctx, "UPDATE Books SET name=COALESCE(NULLIF($1, ''), name), author=COALESCE(NULLIF($2, ''), author), price=COALESCE(NULLIF($3, 0), price) WHERE id=$4", book.Name, book.Author, book.Price, book.ID)
	if err != nil {
		r.logger.Error("error in Update", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("error in Update", err)
		return err
	}
	if rows == 0 {
		r.logger.Error("error in Update", ErrNoChange)
		return ErrNoChange
	}
	return nil
}

func (r *BooksRepo) DeleteByID(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM Books WHERE id=$1", id)
	if err != nil {
		r.logger.Error("error in DeleteByID", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("error in DeleteByID", err)
		return err
	}
	if rows == 0 {
		r.logger.Error("error in DeleteByID", ErrNoChange)
		return ErrNoChange
	}
	return nil
}
