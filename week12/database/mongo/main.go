package main

import (
	"context"
	"fmt"
	"go-course/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Student struct {
	Name  string
	City  string
	Score float32
}

func create(ctx context.Context, collection *mongo.Collection) {
	//插入一个doc
	doc := Student{Name: "张三", City: "北京", Score: 39}
	res, err := collection.InsertOne(ctx, doc)
	database.CheckError(err)
	fmt.Printf("insert id %v\n", res.InsertedID) //每个doc都会有一个全世界唯一的ID(时间+空间唯一)

	//插入多个docs
	docs := []interface{}{Student{Name: "李四", City: "北京", Score: 24}, Student{Name: "王五", City: "南京", Score: 21}}
	manyRes, err := collection.InsertMany(ctx, docs)
	database.CheckError(err)
	fmt.Printf("insert many ids %v\n", manyRes.InsertedIDs)
}

func update(ctx context.Context, collection *mongo.Collection) {
	filter := bson.D{{"city", "北京"}}
	update := bson.D{{"$inc", bson.D{{"score", 5}}}}
	res, err := collection.UpdateMany(ctx, filter, update) //或用UpdateOne
	database.CheckError(err)
	fmt.Printf("update %d doc\n", res.ModifiedCount)
}

func delete(ctx context.Context, collection *mongo.Collection) {
	filter := bson.D{{"name", "张三"}}
	res, err := collection.DeleteMany(ctx, filter) //或用DeleteOne
	database.CheckError(err)
	fmt.Printf("delete %d doc\n", res.DeletedCount)
}

func query(ctx context.Context, collection *mongo.Collection) {
	sort := bson.D{{"name", 1}}                     //1升序，-1降序
	filter := bson.D{{"score", bson.D{{"$gt", 3}}}} //score>3，gt代表greater than
	findOption := options.Find()
	findOption.SetSort(sort)
	findOption.SetLimit(10) //最多返回10个
	findOption.SetSkip(3)   //跳过前3个

	cursor, err := collection.Find(ctx, filter, findOption)
	database.CheckError(err)
	defer cursor.Close(ctx) //关闭迭代器
	for cursor.Next(ctx) {
		var doc Student
		err := cursor.Decode(&doc)
		database.CheckError(err)
		fmt.Printf("%s %s %.2f\n", doc.Name, doc.City, doc.Score)
	}
}

func main() {
	ctx := context.Background()
	option := options.Client().ApplyURI("mongodb://127.0.0.1:27017").
		SetConnectTimeout(time.Second). //连接超时时长
		//AuthSource代表Database
		SetAuth(options.Credential{Username: "tester", Password: "123456", AuthSource: "test"})
	client, err := mongo.Connect(ctx, option)
	database.CheckError(err)
	err = client.Ping(ctx, nil) //Connect没有返回error并不代表连接成功，ping成功才代表连接成功
	database.CheckError(err)
	defer client.Disconnect(ctx) //释放链接

	collection := client.Database("test").Collection("student")
	// create(ctx, collection)
	// update(ctx, collection)
	// delete(ctx, collection)
	query(ctx, collection)
}
