package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"test_task/internal/models"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *models.Subscription) (int, error)
	GetAll(ctx context.Context) ([]models.Subscription, error)
	GetByID(ctx context.Context, id int) (*models.Subscription, error)
	Update(ctx context.Context, id int, input *models.UpdateSubscriptionInput) error
	Delete(ctx context.Context, id int) error
	SumByFilter(ctx context.Context, filter models.SubscriptionFilter) (int, error)
}

type Repository struct {
	cfg     *Config
	db      *pgxpool.Pool
	logger  Logger
	SubRepo SubscriptionRepository
}

func NewRepository(cfg *Config, logger Logger) *Repository {
	return &Repository{
		cfg:    cfg,
		logger: logger,
	}
}

func (r *Repository) Run(_ context.Context) {
}

func (r *Repository) Stop() {
	if r.db != nil {
		r.db.Close()
		r.logger.Info("Database connection closed.")
	}
}

func (r *Repository) Init() error {
	var err error
	r.db, err = newPostgresDB(r.cfg)
	if err != nil {
		return err
	}

	r.SubRepo = NewSubscriptionPostgres(r.db)

	return nil
}
