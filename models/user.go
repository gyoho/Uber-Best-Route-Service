package models

type (
    // User represents the structure of our resource
    // need capitalized to expoert the variables
    // apply alias to use lower case names when delivering json.
    User struct {
        Id          int    `json:"id"`
        Name        string `json:"name"`
        Address     string `json:"address"`
        City        string `json:"city"`
        State       string `json:"state"`
        Zip         string `json:"zip"`
        Coordinate  Coord  `json:"coordinate"`
    }

    Coord struct {
        Lan float64 `json:"lan"`
        Lng float64 `json:"lng"`
    }
)
