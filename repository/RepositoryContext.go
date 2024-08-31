package repository

import (
	"context"
	"gorm.io/gorm"
	"log/slog"
)

// Repository Handling Database transactions cleaner and faster
type Repository interface {
	Create(ctx context.Context, entity interface{}) error
	Query(ctx context.Context, entity interface{}) error
	Transaction(ctx context.Context, fn func(repo Repository) error) error
}

type RepoService struct {
	Db *gorm.DB // implements the repository interface and has the current db conn as a state
}

func (r *RepoService) Create(ctx context.Context, entity interface{}) error {
	slog.Info("Create", "req_id", ctx.Value("req_id"))
	return r.Db.Create(entity).Error
}

func (r *RepoService) Query(ctx context.Context, entity interface{}) error {
	slog.Info("Query", "req_id", ctx.Value("req_id"))
	return r.Db.Where("short_url = ?", ctx.Value("hash")).First(entity).Error
}

func (r *RepoService) RepoTxn(tx *gorm.DB) Repository {
	// create a new repo with the given transaction
	return &RepoService{
		Db: tx,
	}
}

func (r *RepoService) Transaction(ctx context.Context, fn func(repo Repository) error) error {
	slog.Info("Transaction for", "req_id", ctx.Value("req_id"))
	tx := r.Db.Begin() // beginning a database transaction

	if err := tx.Error; err != nil {
		return tx.Error
	}

	repo := r.RepoTxn(tx) // creating a new repo with the existing transaction

	err := fn(repo) // Performing requested DB transaction
	if err != nil {
		tx.Rollback() // Rollback rollbacks the changes in a transaction
		return err
	}

	return tx.Commit().Error // Commit commits the changes in a transaction
}
