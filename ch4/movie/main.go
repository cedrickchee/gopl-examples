// Movie prints Movies as JSON.
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"`
	Actors []string
}

var movies = []Movie{
	{Title: "Casablance", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
	// ...
}

func main() {
	{
		// Marshal
		data, err := json.Marshal(movies)
		if err != nil {
			log.Fatalf("JSON marshalling failed: %s", err)
		}
		fmt.Printf("%s\n", data)
	}

	{
		// MarshalIndent
		//
		// This compact representation contains all the information but it’s hard to
		// read. For human consumption, a variant called `json.MarshalIndent`
		// produces neatly indented output.
		data, err := json.MarshalIndent(movies, "", "    ")
		if err != nil {
			log.Fatalf("JSON marshalling failed: %s", err)
		}
		fmt.Printf("%s\n", data)

		// Unmarshal
		var titles []struct{ Title string }
		if err := json.Unmarshal(data, &titles); err != nil {
			log.Fatalf("JSON unmarshaling failed: %s", err)
		}
		fmt.Println(titles) // "[{Casablanca} {Cool Hand Luke} {Bullitt}]"
	}
}

/*
Output:

[{"Title":"Casablance","released":1942,"Actors":["Humphrey Bogart","Ingrid
Bergman"]},{"Title":"Cool Hand Luke","released":1967,"color":true,
"Actors":["Paul Newman"]},{"Title":"Bullitt","released":1968,"color":true,
"Actors":["Steve McQueen","Jacqueline Bisset"]}]

[
    {
        "Title": "Casablance",
        "released": 1942,
        "Actors": [
            "Humphrey Bogart",
            "Ingrid Bergman"
        ]
    },
    {
        "Title": "Cool Hand Luke",
        "released": 1967,
        "color": true,
        "Actors": [
            "Paul Newman"
        ]
    },
    {
        "Title": "Bullitt",
        "released": 1968,
        "color": true,
        "Actors": [
            "Steve McQueen",
            "Jacqueline Bisset"
        ]
    }
]
*/
