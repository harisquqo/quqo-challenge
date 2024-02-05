package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)


func InitRouter(p *base.Persistence) *gin.Engine {

	r := gin.Default()
	p.Automigrate()
	ProductRoutes(r, p)
	InventoryRoutes(r, p)
	WarehouseRoutes(r, p)
	ImageRoutes(r, p)
	CategoryRoutes(r, p)
	CustomerRoutes(r, p)
	OrderRoutes(r, p)
	return r
}