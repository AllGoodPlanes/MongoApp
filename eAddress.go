package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"os"
)

const (
	MongoDBHosts = "xxxxx"
	AuthDatabase = "xxxxx"
	AuthUserName = "xxxxx"
	AuthPassword = "xxxxx"
	TestDatabase = "xxxxx"
)

type AddressData struct {
	Name  string `bson:"Name"`
	Email string `bson:"Email"`
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/display", display)
	fmt.Println("listening...")
	err := http.ListenAndServe(GetPort(), nil)
	if err != nil {
		panic(err)
	}
}

// Get the Port from the environment
func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func root(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, rootForm)
}

const rootForm = `<!doctype html>
<html>
    <head>
<style>
  body {background-color:blue}
  h1   {color:orange}
  h2   {color: yellow}
  p    {color: white}
legend {color: yellow}
</style>
        <title>TITLE - Fetch Some Data</title>
    </head>
    <body>
        <h1>Description - A simple app to fetch data from a cloud based document database</h1>
        <h2>Form Demo to retrieve e.mail address</h2>
        <form>
              <fieldset>
                        <legend>Enter the search criteria here, e.g. "Maxine":</legend>
                        <p>
                        <label>Name</label>
                       <form action="/display" method="post" accept-charset="utf-8" class="pure-form">
                              <input type="text" name="name" placeholder="name" />
                              <input type="submit" value=".. and query database!" formaction="/display"/>
                        </form>
                        </p>
              </fielsset>
        </form>
    </body>
</html>`

var displayTemplate = template.Must(template.New("display").Parse(displayTemplateHTML))

func display(w http.ResponseWriter, r *http.Request) {

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		fmt.Println("error1")
	}

	mongoSession.SetMode(mgo.Monotonic, true)

	collection := mongoSession.DB(TestDatabase).C("AddressData")

	var result AddressData

	err = collection.Find(bson.M{"Name": r.FormValue("name")}).One(&result)

	if result.Email != "" {
		errn := displayTemplate.Execute(w, "The email id you wanted is: "+result.Email)
		if errn != nil {
			http.Error(w, errn.Error(), http.StatusInternalServerError)
		}
	} else {
		displayTemplate.Execute(w, "Sorry... The email id you wanted does not exist.")
	}
}

const displayTemplateHTML = `<!doctype html>
<html>
    <head>
<style>
  body {background-color:blue}
  h1   {color:orange}
  h2   {color: yellow}
  p    {color: white}
legend {color: yellow}
</style>
        <title>TITLE - Fetch Some Data</title>
    </head>
    <body>
        <h1>Description - A simple app to fetch data from a cloud based document database</h1>
        <h2>Form Demo to retrieve e.mail address</h2>
        <form>
              <fieldset>
                        <legend>Enter the search criteria here, e.g. "Maxine":</legend>
                        <p>
                        <p><b>{{html .}}</b></p>
                        <p><a href="/">Start again!</a></p>
                        </p>
              </fielsset>
        </form>
    </body>
</html>`
