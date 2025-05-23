package main

import (
	"log"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestConnectDatabase(t *testing.T) {
	os.Setenv("DB_USERNAME", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_HOST", "test")
	os.Setenv("DB_PORT", "test")
	os.Setenv("DB_DATABASE", "test")

	database, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Test veritabanına bağlanılamadı: %v", err)
	}
	DB = database

	err = DB.AutoMigrate(&Stock{}, &Store{})
	if err != nil {
		t.Fatalf("Tablolar oluşturulamadı: %v", err)
	}
	if DB == nil {
		t.Error("Veritabanı bağlantısı başlatılamadı")
	}
}

func TestCheckStock(t *testing.T) {

	TestConnectDatabase(t)
	DB.Create(&Stock{ProductSku: "SKU-001", StoreId: 1, ProductPiece: 10, SizeID: 101})
	DB.Create(&Stock{ProductSku: "SKU-001", StoreId: 2, ProductPiece: 5, SizeID: 101})
	DB.Create(&Stock{ProductSku: "SKU-001", StoreId: 1, ProductPiece: 3, SizeID: 102})

	stockMap, err := checkStock("SKU-001", 101)
	if err != nil {
		t.Fatalf("checkStock hatası: %v", err)
	}

	if stockMap[1] != 10 {
		t.Errorf("SKU-001, Beden 101, Mağaza 1 için beklenen stok 10, ancak %d bulundu", stockMap[1])
	}
	if stockMap[2] != 5 {
		t.Errorf("SKU-001, Beden 101, Mağaza 2 için beklenen stok 5, ancak %d bulundu", stockMap[2])
	}
	if _, ok := stockMap[3]; ok {
		t.Errorf("SKU-001, Beden 101, Mağaza 3 için beklenmeyen stok bulundu")
	}
}
func TestCheckStock_NoStockFound(t *testing.T) {
	TestConnectDatabase(t)
	stockMap, err := checkStock("NON_EXISTENT_SKU", 999)
	if err != nil {
		t.Fatalf("checkStock beklenmedik hata döndürdü: %v", err)
	}
	if len(stockMap) != 0 {
		t.Errorf("stok yok, ancak %d eleman bulundu", len(stockMap))
	}
}

func TestGetPreferredStores(t *testing.T) {
	TestConnectDatabase(t)

	DB.Create(&Store{ID: 1, StoreName: "Mağaza A", StorePriority: 1, StoreMax: 100})
	DB.Create(&Store{ID: 2, StoreName: "Mağaza B", StorePriority: 2, StoreMax: 150})
	DB.Create(&Store{ID: 3, StoreName: "Mağaza C", StorePriority: 3, StoreMax: 200})

	DB.Create(&Stock{ProductSku: "SKU-002", StoreId: 1, ProductPiece: 5, SizeID: 201})
	DB.Create(&Stock{ProductSku: "SKU-002", StoreId: 2, ProductPiece: 0, SizeID: 201})
	DB.Create(&Stock{ProductSku: "SKU-002", StoreId: 3, ProductPiece: 15, SizeID: 201})

	preferredStores, err := getPreferredStores("SKU-002")
	if err != nil {
		t.Fatalf("getPreferredStores hatası: %v", err)
	}

	if len(preferredStores) != 2 {
		t.Errorf("Beklenen mağaza sayısı 2, ancak %d bulundu", len(preferredStores))
	}
	if preferredStores[0]["store_id"].(uint) != 1 {
		t.Errorf("İlk mağaza Store ID 1 olmalıydı, ancak %d bulundu", preferredStores[0]["store_id"])
	}
	if preferredStores[1]["store_id"].(uint) != 3 {
		t.Errorf("İkinci mağaza Store ID 3 olmalıydı, ancak %d bulundu", preferredStores[1]["store_id"])
	}

	if preferredStores[0]["stock"].(int) != 0 {
		t.Errorf("Mağaza A için başlangıç stok 0 olmalıydı, ancak %d bulundu", preferredStores[0]["stock"])
	}
}
func TestGetPreferredStores_NoStoresExist(t *testing.T) {
	TestConnectDatabase(t)
	preferredStores, err := getPreferredStores("ANY_SKU")
	if err != nil {
		t.Fatalf("getPreferredStores beklenmedik hata döndürdü: %v", err)
	}
	if len(preferredStores) != 0 {
		t.Errorf("Mağaza bulunamadı, ancak %d eleman bulundu", len(preferredStores))
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	log.Println("Testler tamamlandı")
	os.Exit(code)
}
