package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
	"trekyourworld/db"
	"trekyourworld/response"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"oss.nandlabs.io/golly/rest/server"
)

type TrekData struct {
	Org        string   `json:"org" bson:"org"`
	Uuid       string   `json:"uuid" bson:"uuid"`
	Title      string   `json:"title" bson:"title"`
	Url        string   `json:"url" bson:"url"`
	Elevation  string   `json:"elevation" bson:"elevation"`
	Duration   string   `json:"duration" bson:"duration"`
	Cost       string   `json:"cost" bson:"cost"`
	Difficulty []string `json:"difficulty" bson:"difficulty"`
	Location   string   `json:"location" bson:"location"`
	Tags       []string `json:"tags" bson:"tags"`
}

type TrekFilters struct {
	Organiser  []string `json:"organiser"`
	Location   []string `json:"location"`
	Duration   []string `json:"duration"`
	Difficulty []string `json:"difficulty"`
}

func FindAllTreks(ctx server.Context) {
	collection, err := db.GetCollection("treks_information")
	if err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Failed to get database collection")
	}

	contxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{bson.E{Key: "$unwind", Value: "$treks"}},
		{bson.E{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "title", Value: "$treks.title"},
		}}},
	}

	cur, err := collection.Aggregate(contxt, pipeline)
	if err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Failed to fetch data from db")
	}
	defer cur.Close(contxt)

	var titles []string
	for cur.Next(contxt) {
		var result bson.M
		if err := cur.Decode(&result); err != nil {
			logger.ErrorF("Failed to decode trek information:", err)
			continue
		}
		// Append title to the slice if it exists in the document
		if title, ok := result["title"].(string); ok {
			titles = append(titles, title)
		}
	}

	if err := cur.Err(); err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Cursor Error")
	}

	response.JSON(ctx.HttpResWriter(), http.StatusOK, titles)
}

// search trek by name
func SearchTrek(ctx server.Context) {
	trekName := ctx.GetRequest().URL.Query().Get("trekName")
	logger.Info(trekName)
	if trekName == "" {
		findAllTreksData(ctx)
	} else {
		findTrekByName(ctx, trekName)
	}
}

func findAllTreksData(ctx server.Context) {
	collection, err := db.GetCollection("treks_information")
	if err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Failed to get database collection")
	}

	contxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{bson.E{Key: "$unwind", Value: "$treks"}},
		{bson.E{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0}, // setting the value=0 excludes the field
			{Key: "org", Value: 1},
			{Key: "title", Value: "$treks.title"},
			{Key: "uuid", Value: "$treks.uuid"},
			{Key: "url", Value: "$treks.url"},
			{Key: "elevation", Value: "$treks.elevation"},
			{Key: "duration", Value: "$treks.duration"},
			{Key: "cost", Value: "$treks.cost"},
			{Key: "difficulty", Value: "$treks.difficulty"},
			{Key: "location", Value: "$treks.location"},
			{Key: "distance", Value: "$treks.distance"},
		}}},
	}

	cur, err := collection.Aggregate(contxt, pipeline)
	if err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Failed to fetch data from db")
	}
	defer cur.Close(contxt)

	var treks []TrekData
	for cur.Next(contxt) {
		var result TrekData
		if err := cur.Decode(&result); err != nil {
			logger.ErrorF("Failed to decode trek information:", err)
			continue
		}
		treks = append(treks, result)
	}

	if err := cur.Err(); err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Cursor Error")
	}

	response.JSON(ctx.HttpResWriter(), http.StatusOK, treks)
}

func findTrekByName(ctx server.Context, trekName string) {
	collection, err := db.GetCollection("treks_information")
	if err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Failed to get database collection")
	}

	contxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	regexPattern := fmt.Sprintf(".*%s.*", regexp.QuoteMeta(trekName))
	regex := bson.D{{Key: "$regex", Value: regexPattern}}

	pipeline := mongo.Pipeline{
		{bson.E{Key: "$unwind", Value: "$treks"}},
		{bson.E{Key: "$match", Value: bson.D{
			bson.E{Key: "treks.title", Value: regex},
		}}},
		{bson.E{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0}, // setting the value=0 excludes the field
			{Key: "org", Value: 1},
			{Key: "title", Value: "$treks.title"},
			{Key: "uuid", Value: "$treks.uuid"},
			{Key: "url", Value: "$treks.url"},
			{Key: "elevation", Value: "$treks.elevation"},
			{Key: "duration", Value: "$treks.duration"},
			{Key: "cost", Value: "$treks.cost"},
			{Key: "difficulty", Value: "$treks.difficulty"},
			{Key: "location", Value: "$treks.location"},
			{Key: "distance", Value: "$treks.distance"},
		}}},
	}

	cur, err := collection.Aggregate(contxt, pipeline)
	if err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Failed to fetch data from db")
	}
	defer cur.Close(contxt)

	var treks []TrekData
	for cur.Next(contxt) {
		var result TrekData
		if err := cur.Decode(&result); err != nil {
			logger.ErrorF("Failed to decode trek information:", err)
			continue
		}
		treks = append(treks, result)
	}

	if err := cur.Err(); err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Cursor Error")
	}

	response.JSON(ctx.HttpResWriter(), http.StatusOK, treks)

}

func FilterTreks(ctx server.Context) {
	logger.Info("inside iflter treks")
	collection, err := db.GetCollection("treks_information")
	if err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Failed to get database collection")
	}

	contxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var trekFilter TrekFilters
	if err := json.NewDecoder(ctx.GetRequest().Body).Decode(&trekFilter); err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusBadRequest, "invalid request payload")
		return
	}

	generatedQuery := buildFilterQuery(trekFilter)

	// Define the pipeline
	pipeline := mongo.Pipeline{
		{bson.E{Key: "$unwind", Value: "$treks"}},
		{bson.E{Key: "$match", Value: generatedQuery}},
		{bson.E{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0}, // setting the value=0 excludes the field
			{Key: "org", Value: 1},
			{Key: "title", Value: "$treks.title"},
			{Key: "uuid", Value: "$treks.uuid"},
			{Key: "url", Value: "$treks.url"},
			{Key: "elevation", Value: "$treks.elevation"},
			{Key: "duration", Value: "$treks.duration"},
			{Key: "cost", Value: "$treks.cost"},
			{Key: "difficulty", Value: "$treks.difficulty"},
			{Key: "location", Value: "$treks.location"},
			{Key: "distance", Value: "$treks.distance"},
		}}},
	}

	cur, err := collection.Aggregate(contxt, pipeline)
	if err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Failed to fetch data from db")
	}
	defer cur.Close(contxt)

	var treks []TrekData
	for cur.Next(contxt) {
		var result TrekData
		if err := cur.Decode(&result); err != nil {
			logger.ErrorF("Failed to decode trek information:", err)
			continue
		}
		treks = append(treks, result)
	}

	if err := cur.Err(); err != nil {
		response.Error(ctx.HttpResWriter(), http.StatusInternalServerError, "Cursor Error")
	}

	response.JSON(ctx.HttpResWriter(), http.StatusOK, treks)
}

// buildFilterQuery generates a MongoDB filter query based on TrekFilters.
func buildFilterQuery(filterInfo TrekFilters) bson.M {
	query := bson.M{}

	if len(filterInfo.Organiser) > 0 {
		query["org"] = bson.M{"$in": filterInfo.Organiser}
	}

	if len(filterInfo.Location) > 0 {
		query["treks.location"] = bson.M{"$in": filterInfo.Location}
	}

	if len(filterInfo.Duration) > 0 {
		query["treks.duration"] = bson.M{"$in": filterInfo.Duration}
	}

	if len(filterInfo.Difficulty) > 0 {
		query["treks.difficulty"] = bson.M{"$in": filterInfo.Difficulty}
	}

	return query
}
