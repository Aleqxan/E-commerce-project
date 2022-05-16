package tokens

import (
	"log"
	"os"
	"time"
	"context"

	"github.com/Aleqxan/E-commerce-project/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

type Signeddetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

var UserData *mongo.Collection = database.UserData(database.Client, "Users")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email string, firstname string, lastname string, uid string) (signedtoken string, signedrefreshtoken string, err error) {

	claims := &Signeddetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	refrestoken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshclaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return token, refrestoken, err
}

func ValidateToken(signedtoken string) (claims *SignedDetails, msg string){
	token, err := jwt.ParsewihClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token)(interface{}, error){
		return []byte(SECRET_KEY), nil
	})

	if err != nil{
		msg = err.Error()
		return
	}

	claims, ok := token.Claims(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	cliams.ExpiresAt < time.Now().Local().Unix(){
		msg = "token is already expired"
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string) {

	var ctx, cancel := context.WithTimeOut(context.Backgound(), 100*time.Second)
	
	var updateobj primitive.D

	updateobj = append (updateobj,bson.E(Key:"token", Value: signedtoken))
	updateobj = append (updateobj,bson.E(Key:"refresh_token", Value: signedrefreshtoken))
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateobj = append (updateobj,bson.E(Key:"updateat", Value: update_at))

	upsert := true

	filter := bson.M{"user_id": userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value:updateobj},
	}
	&opt)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
	
}
