package repository

import (
	"context"
	"fmt"
	"strings"
	"test_task/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionPostgres struct {
	db *pgxpool.Pool
}

func NewSubscriptionPostgres(db *pgxpool.Pool) *SubscriptionPostgres {
	return &SubscriptionPostgres{db: db}
}

func (r *SubscriptionPostgres) Create(ctx context.Context, sub *models.Subscription) (int, error) {
	query := `
       INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
       VALUES ($1, $2, $3, $4, $5)
       RETURNING id;
    `
	row := r.db.QueryRow(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)
	var subscriptionID int
	if err := row.Scan(&subscriptionID); err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса на создание подписки: %w", err)
	}
	sub.ID = subscriptionID
	return subscriptionID, nil
}

func (r *SubscriptionPostgres) GetAll(ctx context.Context) ([]models.Subscription, error) {
	query := `
       SELECT id, service_name, price, user_id, start_date, end_date, created_at
       FROM subscriptions;
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса на получение всех подписок: %w", err)
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var s models.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки подписки: %w", err)
		}
		subscriptions = append(subscriptions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации строк при получении всех подписок: %w", err)
	}

	return subscriptions, nil
}

func (r *SubscriptionPostgres) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	query := `
       SELECT id, service_name, price, user_id, start_date, end_date, created_at
       FROM subscriptions
       WHERE id = $1;
    `
	row := r.db.QueryRow(ctx, query, id)

	var s models.Subscription
	err := row.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate, &s.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("подписка с ID %d не найдена: %w", id, err)
		}
		return nil, fmt.Errorf("ошибка выполнения запроса на получение подписки по ID: %w", err)
	}
	return &s, nil
}

func (r *SubscriptionPostgres) Update(ctx context.Context, id int, input *models.UpdateSubscriptionInput) error {
	setValues := []string{}
	args := []interface{}{}
	argIdx := 1

	if input.ServiceName != nil {
		setValues = append(setValues, fmt.Sprintf("service_name = $%d", argIdx))
		args = append(args, *input.ServiceName)
		argIdx++
	}
	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price = $%d", argIdx))
		args = append(args, *input.Price)
		argIdx++
	}
	if input.StartDate != nil {
		setValues = append(setValues, fmt.Sprintf("start_date = $%d", argIdx))
		args = append(args, *input.StartDate)
		argIdx++
	}
	if input.EndDate != nil {
		setValues = append(setValues, fmt.Sprintf("end_date = $%d", argIdx))
		args = append(args, *input.EndDate)
		argIdx++
	}

	if len(setValues) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf(`
       UPDATE subscriptions
       SET %s
       WHERE id = $%d
    `, strings.Join(setValues, ", "), argIdx)

	commandTag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса на обновление подписки: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("не найдена подписка с ID %d для обновления", id)
	}

	return nil
}

func (r *SubscriptionPostgres) Delete(ctx context.Context, id int) error {
	commandTag, err := r.db.Exec(ctx, `DELETE FROM subscriptions WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса на удаление подписки: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("не найдена подписка с ID %d для удаления", id)
	}

	return nil
}

func (r *SubscriptionPostgres) SumByFilter(ctx context.Context, filter models.SubscriptionFilter) (int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	query := "SELECT COALESCE(SUM(price), 0) FROM subscriptions"

	if filter.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIdx))
		args = append(args, filter.UserID)
		argIdx++
	}
	if filter.ServiceName != "" {
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", argIdx))
		args = append(args, filter.ServiceName)
		argIdx++
	}
	if !filter.From.IsZero() {
		conditions = append(conditions, fmt.Sprintf("start_date >= $%d", argIdx))
		args = append(args, filter.From)
		argIdx++
	}
	if !filter.To.IsZero() {
		conditions = append(conditions, fmt.Sprintf("start_date <= $%d", argIdx))
		args = append(args, filter.To)
		argIdx++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	row := r.db.QueryRow(ctx, query, args...)
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса на суммирование по фильтру: %w", err)
	}
	return total, nil
}
