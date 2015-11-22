package utils

import (
    "net/http"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "strings"
    "errors"
    "strconv"
    "../models"
)

const (
    server_token string = "mhYzOb0iETrSZW6xDR_zD6jZhgH1_n3_wbxO_bS4"
    access_token string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicHJvZmlsZSIsInJlcXVlc3QiLCJoaXN0b3J5Il0sInN1YiI6ImMxZmNiNzdjLTI5MDEtNDY2Zi05ZTIyLTFlMTExNmFhZWQwNiIsImlzcyI6InViZXItdXMxIiwianRpIjoiOTY1ZTFjNjgtOTY5Ni00ZDQ1LWFlMWQtMTk5MjU5NzQ2ZDJiIiwiZXhwIjoxNDUwNTkzNTM0LCJpYXQiOjE0NDgwMDE1MzMsInVhY3QiOiJ3R1JyMjlLVGRFdUNHY2ExcHJtVFZRQmhjWmlEc3UiLCJuYmYiOjE0NDgwMDE0NDMsImF1ZCI6ImtBMzNWSHBIV1ZJSUotZXIyVVVodmVwNTZLcm5iMGhmIn0.FUfhN28mAG6_xpSShae8wvTsIcXaH6eA19d056YooD8LTfdxm3vkyLTpm8buiAov9sJY3ww-F6xcKRlyNn9vAzN68jieOqZycJH4XDBh3jKP-qTuc__6N0jbTY4LmWmuCj0qk2oT6g7ERooL7JLKWFNggf9qQYyuX5JB9kJWIzbvB2bHr5ZopCEg6x0pLp1dFmvbrxDmx_QcV_poqA18RKrdvHJ-HgKbTIlGFRHGqg6Wjh5hUtOMOL1-JeCJHvc7DrqDNgVA1uo_GDPpO5a-lWSSwEVSET76A8kNu0JO-ewIZSjJh3MfoGa6Fi9cTx1Vk6gyXfYQyvcSuTC0OFCWFg"
)


func GetEstimates(startCoord models.Coord, endCoord models.Coord, estimate *Estimate) error {
    url := "https://sandbox-api.uber.com/v1/estimates/price?start_latitude=" +
            strconv.FormatFloat(startCoord.Lat, 'f', 16, 64) + "&start_longitude=" + strconv.FormatFloat(startCoord.Lng, 'f', 16, 64) +
            "&end_latitude=" + strconv.FormatFloat(endCoord.Lat, 'f', 16, 64) + "&end_longitude=" + strconv.FormatFloat(endCoord.Lng, 'f', 16, 64) +
            "&server_token=" + server_token

    res, err := http.Get(url)
    if err != nil {
		return err
	}

    if !strings.EqualFold(res.Status, "200 OK") {
        return errors.New("Uber Server Error")
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

    // uberX always is the first obj in the array
    uberXEstimates := contents["prices"].([]interface{})[0]
    estimate.Product_id = uberXEstimates.(map[string]interface{})["product_id"].(string)
    estimate.Costs = uberXEstimates.(map[string]interface{})["low_estimate"].(float64)
    estimate.Duration = uberXEstimates.(map[string]interface{})["duration"].(float64)
    estimate.Distance = uberXEstimates.(map[string]interface{})["distance"].(float64)

    return nil
}

func RequestRide(reqBody models.RideRequest) (float64, error) {
    url := "https://sandbox-api.uber.com/v1/requests?access_token=" + access_token
    // form := url.Values{}
    // form.Add("product_id", reqBody.ProductID)
    // form.Add("start_latitude", reqBody.StartLat)
    // form.Add("start_longitude", reqBody.StartLng)
    // form.Add("end_latitude", reqBody.EndLat)
    // form.Add("end_longitude", reqBody.EndLng)

    buf, _ := json.Marshal(reqBody)
    form := bytes.NewBuffer(buf)

    res, err := http.Post(url,"application/json", form)
    if err != nil {
		return 0.0, err
	}

    if !strings.EqualFold(res.Status, "202 Accepted") {
        return 0.0, errors.New("Uber Server Error")
    }

    body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0.0, err
	}
	defer res.Body.Close()

    var contents map[string]interface{}
	err = json.Unmarshal(body, &contents)
    if err != nil {
        return 0.0, err
    }

    // uberX always is the first obj in the array
    eta := contents["eta"].(float64)

    return eta, nil


}
