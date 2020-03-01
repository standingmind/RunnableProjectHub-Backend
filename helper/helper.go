package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"model"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB : This is helper function to connect mongoDB
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func ConnectDB() *mongo.Collection {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("runnableprojecthub").Collection("projects")

	return collection
}

// ErrorResponse : This is error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// GetError : This is helper function to prepare error model.
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func GetError(err error, w http.ResponseWriter) {

	//log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

func GetProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var pattern = ""
	var projectType = ""
	var language = ""

	// var winProject model.WindowProject
	// var androidProject model.AndroidProject
	//check for request parameter
	pattern = r.URL.Query().Get("pattern")
	projectType = r.URL.Query().Get("projectType")
	language = r.URL.Query().Get("language")

	//no filter for default
	var filter bson.M
	if pattern != "" {
		if projectType != "" && language != "" {
			filter = bson.M{"$and": []interface{}{
				bson.M{"projectname": bson.M{"$regex": pattern, "$options": "$i"}},
				bson.M{"language": language},
				bson.M{"projecttype": projectType},
			},
			}
		} else if projectType != "" && language == "" {
			filter = bson.M{"$and": []interface{}{
				bson.M{"projectname": bson.M{"$regex": pattern, "$options": "$i"}},
				bson.M{"projecttype": projectType},
			},
			}
		} else if projectType == "" && language != "" {
			filter = bson.M{"$and": []interface{}{
				bson.M{"projectname": bson.M{"$regex": pattern, "$options": "$i"}},
				bson.M{"language": language},
			},
			}
		} else if projectType == "" && language == "" {
			
			filter = bson.M{"projectname": bson.M{"$regex": pattern, "$options": "$i"}}

		}
	} else if pattern == "" {
		if projectType != "" && language != "" {
			filter = bson.M{"$and": []interface{}{
				bson.M{"language": language},
				bson.M{"projecttype": projectType},
			},
			}
		} else if projectType != "" && language == "" {
			filter = bson.M{"projecttype": projectType}

		} else if projectType == "" && language != "" {
			filter = bson.M{"language": language}

		} else if projectType == "" && language == "" {
			filter = bson.M{}
		}
	}
	// we created Book array
	var projects []interface{}

	//Connection mongoDB with helper class
	collection := ConnectDB()

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), filter)

	if err != nil {
		GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var project model.Project
		// & character returns the memory address of the following variable.
		err := cur.Decode(&project) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}
		// add item our array
		fmt.Println(project)
		projects = append(projects, project)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(projects) // encode similar to serialize process.
}

func CreateProject(w http.ResponseWriter, project interface{}) {
	w.Header().Set("Content-Type", "application/json")

	// we decode our body request params
	// connect db
	collection := ConnectDB()

	// insert our book model.
	result, err := collection.InsertOne(context.TODO(), project)

	if err != nil {
		GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}
