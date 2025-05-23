package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Stock struct {
	ProductSku   string `json:"product_sku" `
	StoreId      uint   `json:"store_id"`
	ProductPiece int    `json:"product_piece"`
	SizeID       int    `json:"size_id"`
}

type Store struct {
	ID            uint   `json:"id"`
	StoreName     string `json:"store_name"`
	StorePriority int    `json:"store_priority"`
	StoreMax      int    `json:"store_max"`
}

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env dosyası yüklenemedi!")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı!")
	}

	DB = database
}

func checkStock(productSku string, sizeId int) (map[uint]int, error) {
	var stockInfo []struct {
		StoreId      uint `json:"store_id"`
		ProductPiece int  `json:"product_piece"`
	}
	err := DB.Table("stocks").
		Where("product_sku = ? AND size_id = ?", productSku, sizeId).
		Select("store_id, SUM(product_piece) as product_piece").
		Group("store_id").
		Scan(&stockInfo).Error

	if err != nil {
		return nil, err
	}

	stockMap := make(map[uint]int)
	for _, stock := range stockInfo {
		stockMap[stock.StoreId] = stock.ProductPiece
	}

	return stockMap, nil
}

func getPreferredStores(productSku string) ([]fiber.Map, error) {
	var stores []struct {
		StoreId       uint   `json:"store_id"`
		StoreName     string `json:"store_name"`
		StorePriority int    `json:"store_priority"`
		StoreMax      int    `json:"store_max"`
	}
	err := DB.Table("stocks").
		Select("stocks.store_id, stores.store_name, stores.store_priority, stores.store_max").
		Joins("JOIN stores ON stocks.store_id = stores.id").
		Where("stocks.product_sku = ? AND stocks.product_piece > 0", productSku).
		Order("stores.store_priority ASC").
		Group("stocks.store_id").
		Scan(&stores).Error

	if err != nil {
		log.Printf("Hata (checkStock): %v", err)
		return nil, err
	}

	var result []fiber.Map
	for _, store := range stores {
		result = append(result, fiber.Map{
			"store_id":       store.StoreId,
			"store_name":     store.StoreName,
			"stock":          0,
			"store_max":      store.StoreMax,
			"store_priority": store.StorePriority,
		})
	}

	return result, nil
}

func main() {
	ConnectDatabase()
	app := fiber.New()

	app.Get("/stock/:sku/:sizeId", func(c *fiber.Ctx) error {
		sku := c.Params("sku")
		sizeIdStr := c.Params("sizeId")
		sizeId := 0
		if sizeIdStr != "" {
			var err error
			sizeId, err = strconv.Atoi(sizeIdStr)
			if err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Geçersiz beden ID'si"})
			}
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "Beden ID'si zorunludur"})
		}

		stockData, err := checkStock(sku, sizeId)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Stok kontrolü başarısız"})
		}

		preferredStores, err := getPreferredStores(sku)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Depo bulunamadı"})
		}

		for i, store := range preferredStores {
			storeId := store["store_id"].(uint)
			if stock, exists := stockData[storeId]; exists {
				preferredStores[i]["stock"] = stock
			} else {
				preferredStores[i]["stock"] = 0
			}
		}

		return c.JSON(fiber.Map{"product_sku": sku, "stores": preferredStores})
	})

	log.Fatal(app.Listen(":3000"))
}
