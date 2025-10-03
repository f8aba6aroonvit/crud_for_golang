package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoInstance เป็นโครงสร้างสำหรับเก็บ Client และ Database Instance
type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

// MI (MongoInstance) เป็นตัวแปร Global ที่ใช้เก็บการเชื่อมต่อ DB
var MI MongoInstance

func ConnectDB() {

	mongoURI := "mongodb://localhost:27018"
	dbName := "mydatabase"

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Error creating new MongoDB client:", err)
	}

	// 3. สร้าง Context สำหรับ Connection (ตั้ง Timeout ที่ 10 วินาที)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // สำคัญ: ปล่อย Context เมื่อฟังก์ชันสิ้นสุด

	// 4. เชื่อมต่อ Client
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// 5. ตรวจสอบการเชื่อมต่อ (Ping)
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	fmt.Println("✅ Successfully connected to MongoDB!")

	MI = MongoInstance{
		Client: client,
		Db:     client.Database(dbName),
	}
}

func GetCollection(collectionName string) *mongo.Collection {
	return MI.Db.Collection(collectionName)
}
