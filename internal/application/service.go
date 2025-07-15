package application

import (
	"context"
	"fmt"
	"strings"
	"test_task/internal/models"
	"test_task/internal/repository"
	"time"

	"github.com/pkg/errors" // Сохраняем pkg/errors как в оригинале
)

// Subscription определяет интерфейс для бизнес-логики, связанной с подписками.
type Subscription interface {
	Create(ctx context.Context, sub *models.Subscription) (int, error)
	List(ctx context.Context) ([]models.Subscription, error)
	GetByID(ctx context.Context, id int) (*models.Subscription, error)
	Update(ctx context.Context, id int, input *models.UpdateSubscriptionInput) error
	Delete(ctx context.Context, id int) error
	SumByPeriod(ctx context.Context, userID, serviceName, fromDateStr, toDateStr string) (int, error)
}

// Service является контейнером для всех сервисов приложения.
type Service struct {
	Subscription
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Subscription: NewSubscriptionService(repo.SubRepo),
	}
}

type SubscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(ctx context.Context, sub *models.Subscription) (int, error) {
	id, err := s.repo.Create(ctx, sub)
	if err != nil {
		return 0, errors.Wrap(err, "SubscriptionService: не удалось создать подписку в репозитории")
	}
	return id, nil
}

func (s *SubscriptionService) List(ctx context.Context) ([]models.Subscription, error) {
	subs, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "SubscriptionService: не удалось получить список всех подписок из репозитория")
	}
	return subs, nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			return nil, fmt.Errorf("SubscriptionService: подписка с ID %d не найдена", id)
		}
		return nil, errors.Wrapf(err, "SubscriptionService: не удалось получить подписку с ID %d из репозитория", id)
	}
	return sub, nil
}

func (s *SubscriptionService) Update(ctx context.Context, id int, input *models.UpdateSubscriptionInput) error {
	err := s.repo.Update(ctx, id, input)
	if err != nil {
		return errors.Wrapf(err, "SubscriptionService: не удалось обновить подписку с ID %d в репозитории", id)
	}
	return nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.Wrapf(err, "SubscriptionService: не удалось удалить подписку с ID %d из репозитория", id)
	}
	return nil
}

func (s *SubscriptionService) SumByPeriod(ctx context.Context, userID, serviceName, fromDateStr, toDateStr string) (int, error) {
	const layout = "2006-01-02"

	fromTime, err := time.Parse(layout, fromDateStr)
	if err != nil {
		return 0, errors.Wrapf(err, "SubscriptionService: неверный формат даты 'от'. Ожидается %s", layout)
	}

	toTime, err := time.Parse(layout, toDateStr)
	if err != nil {
		return 0, errors.Wrapf(err, "SubscriptionService: неверный формат даты 'до'. Ожидается %s", layout)
	}

	if fromTime.After(toTime) {
		return 0, fmt.Errorf("SubscriptionService: дата 'от' (%s) не может быть позже даты 'до' (%s)", fromDateStr, toDateStr)
	}

	filter := models.SubscriptionFilter{
		UserID:      userID,
		ServiceName: serviceName,
		From:        fromTime,
		To:          toTime,
	}

	sum, err := s.repo.SumByFilter(ctx, filter)
	if err != nil {
		return 0, errors.Wrap(err, "SubscriptionService: не удалось рассчитать сумму по фильтру в репозитории")
	}
	return sum, nil
}
