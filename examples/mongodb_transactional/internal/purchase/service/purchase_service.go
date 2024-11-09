package service

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/purchase/repository"
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/transaction/model"
)

type PurchaseService struct {
	purchaseRepository *repository.PurchaseRepository
}

func (p PurchaseService) Purchase(context context.Context, transaction *model.Transaction) error {
	return p.purchaseRepository.PurchaseWithAutomaticTransaction(context, transaction)
}

func NewPurchaseService(
	purchaseRepository *repository.PurchaseRepository,
) *PurchaseService {
	return &PurchaseService{
		purchaseRepository: purchaseRepository,
	}
}
