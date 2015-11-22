package utils

import "../models"

type (
    // Collection type
	Estimates []Estimate

	Estimate struct {
	    Location_id                 string          `json:"location_id"`
	    Coordinate                  models.Coord    `json:"coordinate"`
		Product_id					string			`json:"product_id"`
	    Costs                       float64         `json:"low_estimate"`
	    Duration                    float64         `json:"duration"`
	    Distance                    float64         `json:"distance"`
	}
)

// Implement the sort.Interface
func (slice Estimates) Len() int {
    return len(slice)
}

func (slice Estimates) Less(i, j int) bool {
    if slice[i].Costs < slice[j].Costs {
        return true
    } else if slice[i].Costs == slice[j].Costs{
        return slice[i].Duration < slice[j].Duration
    } else {
        return false
    }
}

func (slice Estimates) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

// ----------------------------------------------

type Trips []models.Trip

// Implement the sort.Interface
func (slice Trips) Len() int {
    return len(slice)
}

func (slice Trips) Less(i, j int) bool {
    if slice[i].Total_uber_costs < slice[j].Total_uber_costs {
        return true
    } else if slice[i].Total_uber_costs == slice[j].Total_uber_costs{
        return slice[i].Total_uber_duration < slice[j].Total_uber_duration
    } else {
        return false
    }
}

func (slice Trips) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}
