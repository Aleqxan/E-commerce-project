package controllers

import (
	"net"
	"net/http"
	"time"

	"github.com/Aleqxan/E-commerce-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc {

	return func(c *gin.Context){
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid code"})
			c.Abort()
			return
		}
		addresses, err := ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		var addresses models.Address

		addresses.Address_id = primitive.NewObjectID()

		if err = c.BindJSON(&addresses); err != nil{
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		
		match_filter := bson.D{{Key:"$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value:bson.D{primitive.E{Key:"path", Value:"addresses"}}}}
		group := bson.D{{Key: "$group", Value:bson.D{primitive.E{Key:"_id", Value:"$address_id"}, {Key:"count", Value: bson,D{primitive.E{Key: "$sum", Value: 1}}}}}}
		pointcursor, err := UserCollection.Aggrgate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.IndentedJSON(500, "Inetrnal server error")
		}

		var addressinfo []bson.M
		pointcursor.All(ctx, &addressinfo); err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressinfo {
		count := address_no["count"]
		size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			c.IndentedJSON(400, "Not allowed")
		}
		defer cancel()
		ctx.Done
	}	
}

func EditHomeAddress() gin.HandlerFunc {

}

func EditWorkAddress() gin.HandlerFunc {

}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", value: bson.D{primitive.E{Key:"address", Value: addrresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "Wrong command")
			return
		}

		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully deleted")
	}
}
