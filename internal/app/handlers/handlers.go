package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danblok/books/internal/app/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BooksRepo interface {
	GetByID(context.Context, string) (*models.Book, error)
	GetAll(context.Context) ([]*models.Book, error)
	Add(context.Context, models.Book) (*models.Book, error)
	Update(context.Context, models.Book) error
	DeleteByID(context.Context, string) error
}

type BooksHandlers struct {
	BooksRepo BooksRepo
	logger    *slog.Logger
}

func New(booksRepo BooksRepo) *BooksHandlers {
	return &BooksHandlers{
		BooksRepo: booksRepo,
		logger:    slog.Default(),
	}
}

func (h *BooksHandlers) RegisterHandlers(path string, gin *gin.Engine) {
	gin.GET(path+"books", h.GetAll)
	gin.GET(path+"book/:id", h.GetByID)
	gin.POST(path+"book/create", h.Add)
	gin.PATCH(path+"book/update", h.Update)
	gin.DELETE(path+"book/:id", h.DeleteByID)
}

func (h *BooksHandlers) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		h.logger.Error("error parsing id", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	book, err := h.BooksRepo.GetByID(ctx.Request.Context(), id)
	if err != nil {
		h.logger.Error("error getting a book", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(200, book)
}

func (h *BooksHandlers) GetAll(ctx *gin.Context) {
	books, err := h.BooksRepo.GetAll(ctx.Request.Context())
	if err != nil {
		h.logger.Error("error getting books", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(200, books)
}

func (h *BooksHandlers) Add(ctx *gin.Context) {
	var newBook models.Book
	err := ctx.Bind(&newBook)
	if err != nil {
		h.logger.Error("error binding json", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	book, err := h.BooksRepo.Add(ctx.Request.Context(), newBook)
	if err != nil {
		h.logger.Error("error adding a new book", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(201, book)
}

func (h *BooksHandlers) Update(ctx *gin.Context) {
	var bookToUpdate models.Book
	err := ctx.Bind(&bookToUpdate)
	if err != nil {
		h.logger.Error("error binding json", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	if _, err := uuid.Parse(bookToUpdate.ID); err != nil {
		h.logger.Error("error parsing id", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	err = h.BooksRepo.Update(ctx.Request.Context(), bookToUpdate)
	if err != nil {
		h.logger.Error("error updating json", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(200)
}

func (h *BooksHandlers) DeleteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		h.logger.Error("error parsing id", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	err := h.BooksRepo.DeleteByID(ctx.Request.Context(), id)
	if err != nil {
		h.logger.Error("error deleting a book", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(200)
}
