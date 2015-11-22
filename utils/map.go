package utils

import (
    "net/http"
    "io/ioutil"
    "encoding/json"
    "strings"
    "errors"
    "../models"
)

func GetCoordinates(usr *models.User) error {
    url := "http://maps.google.com/maps/api/geocode/json?address=" + usr.Address + ", " + usr.City + ", " + usr.State
    url = strings.Replace(url, " ", "+", -1)

    res, err := http.Get(url)
    if err != nil {
		return err
	}

    body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

    var contents map[string]interface{}
	err = json.Unmarshal(body, &contents)
    if err != nil {
        return err
    }

    if !strings.EqualFold(contents["status"].(string), "OK") {
        return errors.New("Coordinates not found")
    }

    results := contents["results"].([]interface{})
    location := results[0].(map[string]interface{})["geometry"].(map[string]interface{})["location"]
    // usr.Coordinate.Lat, err = strconv.ParseFloat(location.(map[string]interface{})["lat"].(string), 64)
    // usr.Coordinate.Lng, err = strconv.ParseFloat(location.(map[string]interface{})["lng"].(string), 64)

    usr.Coordinate.Lat = location.(map[string]interface{})["lat"].(float64)
    usr.Coordinate.Lng = location.(map[string]interface{})["lng"].(float64)

    if err != nil {
        return err
    }

    return nil
}
