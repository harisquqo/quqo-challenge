package base

import (
	"log"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/category_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/customer_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/image_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/inventory_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/order_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/ordereditem_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/product_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/warehouse_entity"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/logger"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base/db"
	"github.com/opensearch-project/opensearch-go"
	"github.com/redis/go-redis/v9"
	storage_go "github.com/supabase-community/storage-go"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// Repositories struct -- Currently only Product Repo
type Persistence struct {
	DB *gorm.DB
	DbRedis *redis.Client
	DbMongo *mongo.Client
	DbOpensearch *opensearch.Client
	DbSupabase *storage_go.Client
	Logger logger.LoggerRepo
}

// Function to create a new repository
func NewPersistence() (*Persistence, error) {
	database, errDatabase := db.NewDB()
	redisDb, errRedisDb := db.NewRedisDB()
	mongoDb, errMongoDb := db.NewMongoDB()
	opensearchDb, errOpensearchDb := db.NewOpenSearchDB()
	supabaseDb, errSupabaseDb := db.NewSupabaseDB()
	logger, errLogger := db.NewLogger()
	if errDatabase != nil {
		log.Fatal(errDatabase)
	}

	if errRedisDb != nil {
		log.Fatal(errRedisDb)
	}

	if errMongoDb != nil {
		log.Fatal(errMongoDb)
	}

	if errOpensearchDb != nil {
		log.Fatal(errOpensearchDb)
	}

	if errSupabaseDb != nil {
		log.Fatal(errSupabaseDb)
	}

	if errSupabaseDb != nil {
		log.Fatal(errLogger)
	}




	return &Persistence{
		DB: database.DB,
		DbRedis: redisDb,
		DbMongo: mongoDb,
		DbOpensearch: opensearchDb,
		DbSupabase: supabaseDb,
		Logger: logger, 
	}, nil

}

//This migrate all tables
func (s *Persistence) Automigrate() error {
	return s.DB.AutoMigrate(&product_entity.Product{}, 
		&inventory_entity.Inventory{},
		&inventory_entity.InventoryLog{}, 
		&warehouse_entity.Warehouse{}, 
		&image_entity.Image{},
		&category_entity.Category{},
		&customer_entity.Customer{},
		&order_entity.Order{},
		&ordereditem_entity.OrderedItem{})
}