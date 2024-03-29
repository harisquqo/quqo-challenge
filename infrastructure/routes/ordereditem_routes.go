package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func OrderedItemRoutes(router *gin.RouterGroup, p *base.Persistence) {
    orderedItems := handlers.NewOrderedItem(p)
    
    router.GET("admin/ordereditems", orderedItems.GetAllOrderedItems)
    router.GET("admin/orders/:order_id/ordereditems", orderedItems.GetAllOrderedItemsForOrder)
}
