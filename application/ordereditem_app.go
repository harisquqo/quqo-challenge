package application

import (
	"fmt"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/ordereditem_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/ordereditem_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/inventories"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/ordereditems"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type OrderedItemApp struct {
	p *base.Persistence
}

func NewOrderedItemApplication(p *base.Persistence) ordereditem_repository.OrderedItemRepository {
	return &OrderedItemApp{p}
}

func (a *OrderedItemApp) ReverseOrder(orderedItems []ordereditem_entity.OrderedItem) map[string]string {
	errorMap := make(map[string]string)

	for _, orderedItem := range orderedItems {
		// Calculate the quantity to add to inventory by reversing the ordered quantity
		quantityToAdd := orderedItem.Quantity

		// Increase inventory for the product
		inventoryRepo := inventories.NewInventoryRepository(a.p)
		err := inventoryRepo.IncreaseInventory(orderedItem.ProductID, quantityToAdd)
		if err != nil {
			errorMap[fmt.Sprintf("product_%d", orderedItem.ProductID)] = err.Error()
		}
	}

	// Return error map, if any
	return errorMap
}

func (a *OrderedItemApp) GetAllOrderedItems() ([]ordereditem_entity.OrderedItem, error) {
	repoOrderedItem := ordereditems.NewOrderedItemsRepository(a.p)
	return repoOrderedItem.GetAllOrderedItems()
}

func (a *OrderedItemApp) GetAllOrderedItemsForOrder(orderId int64) ([]ordereditem_entity.OrderedItem, error) {
	repoOrderedItem := ordereditems.NewOrderedItemsRepository(a.p)
	return repoOrderedItem.GetAllOrderedItemsForOrder(orderId)
}
