package orders

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/order_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/order_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/cache"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"gorm.io/gorm"
)

type OrderRepo struct {
	p *base.Persistence
	c *gin.Context
}

func NewOrderRepository(p *base.Persistence, c *gin.Context) *OrderRepo {
	return &OrderRepo{p, c}
}

var _ order_repository.OrderRepository = &OrderRepo{}

func (o *OrderRepo) SaveOrder(tx *gorm.DB, order *order_entity.Order) (*order_entity.Order, error) {
	span := o.p.Logger.Start(o.c, "implementations/SaveOrder")
	defer span.End()
	if tx == nil {
		var errTx error
		tx := o.p.DB.Begin()
		if tx.Error != nil {
			return nil, errors.New("failed to start transaction")
		}
	
		// Defer rollback in case of panic
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			} else if errTx != nil {
				tx.Rollback()
			} else {
				errC := tx.Commit().Error
				if errC != nil {
					tx.Rollback()
				}
			}
		}()
	}
	err := tx.Debug().Create(&order).Preload("OrderedItems").Error
	if err != nil {
		// tracer.Span.RecordError(err)
		fmt.Println("Failed to create order")
		fmt.Println(err)
		return nil, err
	}
	
	o.p.Logger.Info("implementations/SaveOrder", map[string]interface{}{"order": order})
	return order, nil
}


func (o *OrderRepo) GetOrder(id int64) (*order_entity.Order, error) {
	span := o.p.Logger.Start(o.c, "implementations/GetOrder")
	defer span.End()
	var order *order_entity.Order
	cacheRepo := cache.NewCacheRepository("Redis", o.p)
	_ = cacheRepo.GetKey(fmt.Sprintf("%v_ORDER", id), &order)
	if order == nil {
		err := o.p.DB.Debug().Preload("OrderedItems").Where("id = ?", id).Take(&order).Error
		if err != nil {
			fmt.Println("Failed to get order")
		}
		if order != nil && order.ID > 0 {
			_ = cacheRepo.SetKey(fmt.Sprintf("%v_ORDER", id), order, time.Minute * 15)
		}
	}


	return order, nil
}

func (o *OrderRepo) GetAllOrders() ([]order_entity.Order, error) {
	var orders []order_entity.Order
	err := o.p.DB.Debug().Preload("OrderedItems").Find(&orders).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return orders, nil
}

func (o *OrderRepo) UpdateOrder(order *order_entity.Order) (*order_entity.Order, error) {
	span := o.p.Logger.Start(o.c, "implementations/UpdateOrder")
	defer span.End()
	cacheRepo := cache.NewCacheRepository("Redis", o.p)


	err := o.p.DB.Debug().Where("id = ?", order.ID).Updates(&order).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	_ = cacheRepo.SetKey(fmt.Sprintf("%v_ORDER", order.ID), order, time.Minute * 15)

	return order, nil
}

func (o *OrderRepo) DeleteOrder(id int64) error {
	var order order_entity.Order

	err := o.p.DB.Debug().Where("id = ?", id).Delete(&order).Error
	
	cacheRepo := cache.NewCacheRepository("Redis", o.p)

	cacheRepo.DelKey(fmt.Sprintf("%v_ORDER", id))
	if err != nil {
		return errors.New("database error, please try again")
	}

	return nil
}