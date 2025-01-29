package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Beneficiary struct {
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
	Msisdn    string `json:"msisdn"`
	Country   string `json:"country"`
	Address   string `json:"address"`
}

type Transaction struct {
	ID          uint    `gorm:"primaryKey"`
	Operator    string  `json:"operator"`
	Amount      float64 `json:"amount"`
	Msisdn      string  `json:"msisdn"`
	Country     string  `json:"country"`
	Status      string  `gorm:"default:pending"`
	CreatedAt   time.Time
	Beneficiary Beneficiary `json:"beneficiary" gorm:"embedded"`
}

var db *gorm.DB

func initDB() {
	dsn := "mamlakadev:@Mamlaka2021@tcp(localhost:3306)/westsandbox?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	db.AutoMigrate(&Transaction{})
}

func sendTransaction(c *gin.Context) {
	var txn Transaction
	if err := c.ShouldBindJSON(&txn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txn.Status = "successful" // Mock successful transaction
	db.Create(&txn)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transaction Successful",
		"transaction": txn,
	})
}

func payoutTransaction(c *gin.Context) {
	var txn Transaction
	if err := c.ShouldBindJSON(&txn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txn.Status = "payout successful" // Mock payout success
	db.Create(&txn)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Payout Successful",
		"transaction": txn,
	})
}

func checkTransactionStatus(c *gin.Context) {
	id := c.Param("id")
	var txn Transaction
	if err := db.First(&txn, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactionId": txn.ID,
		"status":        txn.Status,
	})
}

func main() {
	initDB()
	r := gin.Default()
	r.POST("/send-transaction", sendTransaction)
	r.POST("/payout-transaction", payoutTransaction)
	r.GET("/transaction-status/:id", checkTransactionStatus)

	r.Run(":8080")
}
