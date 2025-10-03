package controllers

import (
	"context"
	"log"
	"time"

	"train_golang/configs"
	"train_golang/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(c *fiber.Ctx) error {
	var userCollection *mongo.Collection = configs.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	newID := primitive.NewObjectID()

	// 1. **ขั้นตอนที่จำเป็น**: ต้องมีการสร้าง Document 1 อันใน Collection ชั่วคราว
	//    หรือใช้ Collection ที่คุณมั่นใจว่ามี 1 Document อยู่แล้ว
	//    เพื่อให้มี Input Document สำหรับ Pipeline

	// **ตัวอย่างนี้ใช้เทคนิคการสร้าง Document ภายใน Pipeline**
	// **NOTE**: หาก Collection 'users' มี Document อยู่แล้ว มันจะถูกใช้เป็น Input Document
	//         ซึ่งจะทำให้ Aggregation ทำงานได้ แต่จะทำ $merge ซ้ำตามจำนวน Document
	//         ดังนั้นต้องใส่ $limit: 1 ไว้ก่อน

	pipeline := mongo.Pipeline{
		// 1. $limit: 1 (บังคับให้มี Document อินพุตแค่ 1 ตัว หาก Collection ไม่ว่าง)
		{{"$limit", 1}},

		// 2. $project: สร้าง Document ใหม่โดยใช้ค่าคงที่จาก Request Body
		//    (โดยละเลยค่าจาก Document อินพุตเดิม)
		{{"$project", bson.D{
			{"_id", newID},
			{"name", user.Name},
			{"email", user.Email},
			{"password", user.Password},
			{"_id_old", "$_id"}, // เก็บ ID เดิมไว้ (ถ้ามี)
		}}},

		// 3. $unset: ลบฟิลด์ที่ไม่ต้องการ
		{{"$unset", bson.A{"_id_old"}}},

		// 4. $merge: แทรก Document ใหม่เข้าไปใน Collection 'users'
		{{"$merge", bson.D{
			{"into", "users"},
			{"whenMatched", "fail"},
			{"whenNotMatched", "insert"},
		}}},
	}

	// 4. รัน Aggregate
	// ถ้า Collection 'users' ว่างเปล่า (ไม่มี Document) ไพพ์ไลน์จะไม่ทำงาน
	cursor, err := userCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Aggregate Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error creating user with pipeline", "error": err.Error()})
	}

	// ต้อง Close Cursor เพื่อให้ $merge ทำงาน
	cursor.Close(ctx)

	// 5. ส่ง Response กลับ
	// ยืนยันการสร้างด้วย ID
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully with aggregation", "userId": newID.Hex()})
}
