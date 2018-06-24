package controllers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/gin-gonic/gin/binding"
	"github.com/swhite24/go-lambda-test/pkg/models"

	"github.com/gin-gonic/gin"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type (
	// UserController provides methods for interacting with users
	UserController struct {
		tname string
		db    *dynamodb.DynamoDB
	}
)

// ServeUsers wires the user controller into the provided engine
func ServeUsers(e *gin.Engine, tname string) {
	// Setup db
	config := aws.NewConfig().WithRegion("us-east-1")
	s, _ := session.NewSession()
	db := dynamodb.New(s, config)

	uc := &UserController{tname, db}
	e.GET("/user", uc.listUsers)
	e.GET("/user/:id", uc.getUser)
	e.POST("/user", uc.createUser)
}

func (uc *UserController) createUser(c *gin.Context) {
	var err error
	var user models.User
	var av map[string]*dynamodb.AttributeValue
	if err = binding.JSON.Bind(c.Request, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "unable to parse user"})
		return
	}

	user.ID = genHash(16)
	if av, err = dynamodbattribute.MarshalMap(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "unable to decode user into attribute values"})
		return
	}

	_, err = uc.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(uc.tname),
		Item:      av,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "unable to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "user": user})
}

func (uc *UserController) listUsers(c *gin.Context) {
	var err error
	var res *dynamodb.ScanOutput
	var users []*models.User

	res, err = uc.db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(uc.tname),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "unable to list users", "err": err.Error()})
		return
	}

	dynamodbattribute.UnmarshalListOfMaps(res.Items, &users)
	c.JSON(http.StatusOK, gin.H{"success": true, "users": users})
}

func (uc *UserController) getUser(c *gin.Context) {
	var err error
	var res *dynamodb.GetItemOutput
	var user models.User
	id := c.Param("id")

	res, err = uc.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(uc.tname),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "unable to get user", "err": err.Error()})
		return
	}
	dynamodbattribute.UnmarshalMap(res.Item, &user)
	c.JSON(http.StatusOK, gin.H{"success": true, "user": user})
}

// genHash creates a random hash of alphanumeric characters of provided length
func genHash(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
