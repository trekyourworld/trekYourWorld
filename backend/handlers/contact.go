package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"trekyourworld/db"
	"trekyourworld/response"

	"oss.nandlabs.io/golly/rest/server"
)

type ContactUsSchema struct {
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Email     string `json:"email" bson:"email"`
	Message   string `json:"message" bson:"message"`
}

func ContactUs(ctx server.Context) {
	var contactUs ContactUsSchema
	if err := json.NewDecoder(ctx.GetRequest().Body).Decode(&contactUs); err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusBadRequest, "invalid request payload")
		return
	}

	collection, err := db.GetCollection("contact_us")
	if err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Failed to get database collection")
	}

	contxt := context.Background()

	result, err := collection.InsertOne(contxt, contactUs)
	if err != nil {
		logger.ErrorF("Failed to insert document: %v", err)
		response.Error(ctx.HttpResWriter(), http.StatusOK, "error adding feedback")
	}

	logger.InfoF("Inserted document with ID: %v", result.InsertedID)
	response.JSON(ctx.HttpResWriter(), http.StatusOK, "feedback added successfully")
}
