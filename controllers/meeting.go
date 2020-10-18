package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MsgResponse represents JSON Response structure
type MsgResponse struct {
	Status  int
	Message string
	Data    []Meeting
	Time    time.Time
}

// Meeting represents Meeting collection structure
type Meeting struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Participants string    `json:"participants"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	CreatedAt    time.Time `json:"created_at"`
}

// DATABASE INSTANCE
var collection *mongo.Collection

// MeetingCollection initialises new collection
func MeetingCollection(c *mongo.Database) {
	collection = c.Collection("meetings")
}

// GetAllMeetings returns all meetings scheduled between given start and end time
func GetAllMeetings(startTime string, endTime string, r *http.Request) MsgResponse {
	meetings := []Meeting{}
	cursor, err := collection.Find(context.TODO(), bson.M{"start_time": startTime, "end_time": endTime})

	if err != nil {
		log.Printf("[x] Error while getting all meetings, Reason: %v\n", err)
		message := MsgResponse{
			http.StatusInternalServerError,
			"Something went wrong",
			make([]Meeting, 0),
			time.Now().UTC(),
		}
		return message
	}

	// Iterate through the returned cursor.
	for cursor.Next(context.TODO()) {
		var meeting Meeting
		cursor.Decode(&meeting)
		meetings = append(meetings, meeting)
	}

	message := MsgResponse{
		http.StatusOK,
		"All Meetings",
		meetings,
		time.Now().UTC(),
	}
	return message
}

// func CreateTodo(c *gin.Context) {
// 	var todo Todo
// 	c.BindJSON(&todo)
// 	title := todo.Title
// 	body := todo.Body
// 	completed := todo.Completed
// 	id := guuid.New().String()

// 	newTodo := Todo{
// 		ID:        id,
// 		Title:     title,
// 		Body:      body,
// 		Completed: completed,
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	_, err := collection.InsertOne(context.TODO(), newTodo)

// 	if err != nil {
// 		log.Printf("Error while inserting new todo into db, Reason: %v\n", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  http.StatusInternalServerError,
// 			"message": "Something went wrong",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"status":  http.StatusCreated,
// 		"message": "Todo created Successfully",
// 	})
// 	return
// }

// func GetSingleTodo(c *gin.Context) {
// 	todoId := c.Param("todoId")

// 	todo := Todo{}
// 	err := collection.FindOne(context.TODO(), bson.M{"id": todoId}).Decode(&todo)
// 	if err != nil {
// 		log.Printf("Error while getting a single todo, Reason: %v\n", err)
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"status":  http.StatusNotFound,
// 			"message": "Todo not found",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  http.StatusOK,
// 		"message": "Single Todo",
// 		"data":    todo,
// 	})
// 	return
// }

// func EditTodo(c *gin.Context) {
// 	todoId := c.Param("todoId")
// 	var todo Todo
// 	c.BindJSON(&todo)
// 	completed := todo.Completed

// 	newData := bson.M{
// 		"$set": bson.M{
// 			"completed":  completed,
// 			"updated_at": time.Now(),
// 		},
// 	}

// 	_, err := collection.UpdateOne(context.TODO(), bson.M{"id": todoId}, newData)
// 	if err != nil {
// 		log.Printf("Error, Reason: %v\n", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  500,
// 			"message": "Something went wrong",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  200,
// 		"message": "Todo Edited Successfully",
// 	})
// 	return
// }

// func DeleteTodo(c *gin.Context) {
// 	todoId := c.Param("todoId")

// 	_, err := collection.DeleteOne(context.TODO(), bson.M{"id": todoId})
// 	if err != nil {
// 		log.Printf("Error while deleting a single todo, Reason: %v\n", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  http.StatusInternalServerError,
// 			"message": "Something went wrong",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  http.StatusOK,
// 		"message": "Todo deleted successfully",
// 	})
// 	return
// }
