package handlers

import (
	"context"
	"net/http"
	"trekyourworld/db"
	"trekyourworld/response"

	"go.mongodb.org/mongo-driver/bson"
	"oss.nandlabs.io/golly/rest/server"
)

type Organisation struct {
	ID    string `json:"id" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Label string `json:"label" bson:"label"`
}

func FindAllOrganisations(ctx server.Context) {
	collection, err := db.GetCollection("organisations")
	if err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Failed to get database collection")
	}

	filter := bson.D{}

	contxt := context.Background()

	cur, err := collection.Find(contxt, filter)
	if err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Failed to fetch data from db")
	}
	defer cur.Close(contxt)

	var orgs []Organisation
	for cur.Next(contxt) {
		var org Organisation
		if err := cur.Decode(&org); err != nil {
			logger.ErrorF("Failed to decode organisation:", err)
			continue
		}
		orgs = append(orgs, org)
	}

	if err := cur.Err(); err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Cursor Error")
	}

	response.JSON(ctx.HttpResWriter(), http.StatusOK, orgs)
}
