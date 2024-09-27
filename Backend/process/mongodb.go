package process

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func data_test() {
	// 创建 MongoDB 客户端连接选项
	clientOptions := options.Client().ApplyURI("mongodb://sbzz_admin:sbzz123@127.0.0.1:27017/admin")
	fmt.Println("Successfully initialize Client")

	// 连接到 MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	// 确保成功连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Could not ping MongoDB:", err)
	}
	fmt.Println("Successfully connected to MongoDB!")

	// 选择数据库和集合
	collection := client.Database("SIDT_FP").Collection("part1")

	// 查询条件，假设查询所有名称为 "John" 的文档
	filter := bson.D{{}}

	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error finding documents:", err)
	}

	// 遍历查询结果
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal("Error decoding document:", err)
		}
		fmt.Println("Found document:", result)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal("Error iterating through cursor:", err)
	}

	// 关闭连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal("Error disconnecting from MongoDB:", err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
