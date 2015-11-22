package utils

import (
    "../models"
    "sort"
)

type (
	// To parse the request body
	Location struct {
	    Starting_from_location_id   string          `json:"starting_from_location_id"`
	    Location_ids                []string        `json:"location_ids"`
	}

	List []Element
	Element struct {
		Ids []string
	}

	ListEstimates []Estimates
)

func FindBestRouteBreadthFirst(allPossibleEstimates ListEstimates) (models.Trip, error) {
    var allTrips Trips

    for _, estimates := range allPossibleEstimates {
        trip := models.Trip{}
        for i := 0; i < len(estimates) - 1; i++ {
            err := GetEstimates(estimates[i].Coordinate, estimates[i+1].Coordinate, &estimates[i+1])
            if err != nil {
                return trip, err
            }

            trip.Best_route_location_ids = append(trip.Best_route_location_ids, estimates[i+1].Location_id)
            trip.Product_ids = append(trip.Product_ids, estimates[i+1].Product_id)
            trip.Total_uber_costs += estimates[i+1].Costs
            trip.Total_uber_duration += estimates[i+1].Duration
            trip.Total_distance += estimates[i+1].Distance
        }
        allTrips = append(allTrips, trip)
    }

    // sort by costs then duration
    sort.Sort(allTrips)
    // return the resonable one
    return allTrips[0], nil
}

func FindBestRouteDepthFirst(estimatesArray Estimates) (models.Trip, Estimates, error) {
    trip := models.Trip{}
    for len(estimatesArray) > 1 {
        // set the start location to the arrival point
        startLoc := estimatesArray[0]
        // remove the arrival point from the destination points
        estimatesArray = estimatesArray[1:]

        var err error
        estimatesArray, err = FindRouteDepthFirst(startLoc, estimatesArray, &trip)
        if err != nil {
            return models.Trip{}, Estimates{}, err
        }
    }

    return trip, estimatesArray, nil
}

func FindRouteDepthFirst(startLoc Estimate, estimatesArray Estimates, trip *models.Trip) (Estimates, error) {
    for i, estimate := range estimatesArray {
        err := GetEstimates(startLoc.Coordinate, estimate.Coordinate, &estimate)
        if err != nil {
            return estimatesArray, err
        }
        estimatesArray[i] = estimate
    }

    // sort by costs then duration
    sort.Sort(estimatesArray)
    trip.Best_route_location_ids = append(trip.Best_route_location_ids, estimatesArray[0].Location_id)
    trip.Product_ids = append(trip.Product_ids, estimatesArray[0].Product_id)
    trip.Total_uber_costs += estimatesArray[0].Costs
    trip.Total_uber_duration += estimatesArray[0].Duration
    trip.Total_distance += estimatesArray[0].Distance

    return estimatesArray, nil
}

func FindAllPossibleRoute(Location_ids []string) List {
    prefix := []string{}
    permutation := List{}
    CalculatePermutation(prefix, Location_ids, &permutation)
    return permutation
}
