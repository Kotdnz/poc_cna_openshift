// EPAM: POC CNA in OpenShift with ConfigMaps
// September, 2020
// Kostiantyn_Nikonenko@epam.com

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const appVer = "v0.1"

// CNAusers is our users in DB
/*
type CNAusers struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // omitempty to protect against zeroed _id insertion
	firstname string             `bson:"string firstname" json:"string firstname"`
	lastname  string             `bson:"string lastname" json:"string lastname"`
	phone     string             `bson:"string phone" json:"string phone"`
	email     string             `bson:"string email" json:"string email"`
} */

// Config predefinition
var confPath string

// Config struct
type Config struct {
	BindAddr string `tool:"bind_addr"`
	WhoAmI   string `tool:"whoami"`
	MongoURI string `tool:"MongoURI"`
	DBuser   string `tool:"mongodb"`
	DBpass   string `tool:"mongodh"`
	DBtable  string `tool:"sampledb"`
}

// NewConfig is for generating new one with default values
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		WhoAmI:   "Client",
		MongoURI: "mongodb://mongodb:mongodb@localhost:27017/sampledb",
		DBuser:   "mongodb",
		DBpass:   "mongodb",
		DBtable:  "sampledb",
	}
}

// main function
func main() {

	// reading config in toml format
	flag.StringVar(&confPath, "config-path", "configs/cna-config.toml", "Path to the config file")
	flag.Parse()
	// reading default values
	config := NewConfig()
	// override by values in config file
	_, err := toml.DecodeFile(confPath, config)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the MongoDB
	// Set client options
	credential := options.Credential{
		Username: config.DBuser,
		Password: config.DBpass,
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(config.MongoURI).SetAuth(credential)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Can't connect:", err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Can't ping DB: ", err)
	} else {
		log.Printf("Connected to MongoDB!")
	}

	// Get a handle for your collection
	collection := client.Database(config.DBtable).Collection("cnausers")

	// handling the end-points
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// status code 200
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<p> Wrong place </p>"+
			"<br> navigate to the url: http://hostname:port/api </br>"+
			"<b>I'm %s</b>", config.WhoAmI)
	})

	// main worker function
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "BackEnd is a %s", config.WhoAmI)
	})

	// main worker function
	http.HandleFunc("/api/new", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)

		if err := req.ParseForm(); err != nil {
			// handle error
			log.Fatal(err)
		}

		newUser := bson.M{
			"firstname": req.FormValue("firstname"),
			"lastname":  req.FormValue("lastname"),
			"phone":     req.FormValue("phone"),
			"email":     req.FormValue("email"),
		}

		_, err := collection.InsertOne(context.Background(), newUser)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(newUser)
	})

	// main worker function
	http.HandleFunc("/api/get", func(w http.ResponseWriter, r *http.Request) {

		cursor, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			log.Fatal("Find error: ", err)
		}

		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			log.Fatal("Find_All error", err)
		}

		jsonResp, merr := json.Marshal(results)
		if merr != nil {
			log.Fatal("Marshaling error: ", err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
		log.Println("Request to api/get has been handled with success.")
	})

	// main worker function
	http.HandleFunc("/api/del", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		log.Println(r)
	})

	// check if our app able to handle the requests
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok!"))
	})

	// Check if our app is live
	http.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Live!"))
	})

	log.Printf("listening on %s", config.BindAddr)

	errr := http.ListenAndServe(config.BindAddr, nil)
	if errr != nil {
		log.Fatalf("Failed to start BackEnd server:%v", errr)
	}
}
