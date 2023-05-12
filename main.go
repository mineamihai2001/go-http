package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	. "github.com/mineamihai2001/cc/tema_1/core"
	db "github.com/mineamihai2001/cc/tema_1/database"
	"go.mongodb.org/mongo-driver/bson"
)

type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	err := godotenv.Load()
	Check(err, "Error loading .env file")
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	server := NewServer(port)
	client := db.NewClient("random_facts")

	// Read All
	server.Get("/facts", func(req *Request, res *Response) {
		results := db.Find[db.Fact](client, "facts", bson.D{})

		res.Json(*results).Status(http.StatusOK)
	})

	// Read with category
	server.Get("/facts/:category", func(req *Request, res *Response) {
		args := req.Params()

		category := db.Category(args[1])
		filter := bson.D{{"category", category}}
		results := db.Find[db.Fact](client, "facts", filter)

		if len(*results) == 0 {
			res.Raw("not found").Status(http.StatusNotFound)
			return
		}

		res.Json(results).Status(http.StatusOK)
	})

	// Create
	server.Post("/facts", func(req *Request, res *Response) {
		var dto db.CreateFactDto
		Body(req, &dto)
		fmt.Println(dto)

		if dto.Category == "" {
			dto.Category = db.GENERAL
		}

		id := db.Insert(client, "facts", dto)
		res.Json(db.CreateFactResponseDto{
			CreatedId: id,
		}).Status(http.StatusCreated)
	})

	// Update
	server.Put("/facts", func(req *Request, res *Response) {
		var dto db.UpdateFactDto
		Body(req, &dto)

		if dto.Category == "" {
			dto.Category = db.GENERAL
		}
		update := bson.M{
			"$set": bson.M{
				"category": dto.Category,
				"question": dto.Question,
				"answer":   dto.Answer,
			},
		}
		result := db.Update(client, "facts", dto.Id, update)
		res.Json(db.UpdateFactResponseDto{
			UpdatedCount: result,
		}).Status(http.StatusOK)
	})

	// Update (raw response)
	server.Put("/facts/raw", func(req *Request, res *Response) {
		var dto db.UpdateFactDto
		Body(req, &dto)

		if dto.Category == "" {
			dto.Category = db.GENERAL
		}
		update := bson.M{
			"$set": bson.M{
				"category": dto.Category,
				"question": dto.Question,
				"answer":   dto.Answer,
			},
		}
		result := db.Update(client, "facts", dto.Id, update)
		res.Raw(fmt.Sprint(result)).Status(http.StatusOK)
	})

	// Delete
	server.Delete("/facts/:id", func(req *Request, res *Response) {
		args := req.Params()

		id := args[1]
		result := db.DeleteOne(client, "facts", id)

		status := http.StatusOK
		if result == 0 {
			status = http.StatusNotFound
		}

		dto := db.DeleteFactResponseDto{
			DeletedCount: result,
		}
		res.Json(dto).Status(status)
	})

	// Delete category
	server.Delete("/facts", func(req *Request, res *Response) {
		var dto db.DeleteFactDto
		Body(req, &dto)

		filter := bson.D{{"category", dto.Category}}
		result := db.DeleteMany(client, "facts", filter)

		status := http.StatusOK
		if result == 0 {
			status = http.StatusNotFound
		}

		resp := db.DeleteFactResponseDto{
			DeletedCount: result,
		}
		res.Json(resp).Status(status)
	})

	server.Run(func() {
		fmt.Println("Server listening on port:", port)
	})
}
