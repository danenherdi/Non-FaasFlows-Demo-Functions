package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Input is the input to the function handler function in the template project
func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		defer r.Body.Close()
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("body is empty"))
		return
	}

	dataBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %s", err.Error())))
		return
	}

	var input Input
	err = json.Unmarshal(dataBytes, &input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %s", err.Error())))
		return
	}

	// Check if the user_id is provided
	if input.UserID == nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("user id is empty"))
		return
	}

	// Request to last ride function to get the passenger's last ride
	reqBody, err := json.Marshal(map[string]interface{}{
		"user_id": *input.UserID,
	})
	req, err := http.NewRequest(
		"POST",
		"http://127.0.0.1:8080/function/last-ride-nonflow",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		fmt.Printf("error in creating new request of function last-ride-nonflow: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error in doing the request of function last-ride-nonflow: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("failed request of last-ride-nonflow")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("failed request of last-ride-nonflow"))
		return
	}

	// Parse the response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	var lastRide Ride
	err = json.Unmarshal(data, &lastRide)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("%+v\n", lastRide)

	// Generate ride history with the last ride repeated 3 times
	rideHistory := RideHistory{
		Rides: []RideSummary{
			{
				Time:        lastRide.Time,
				Destination: lastRide.Destination,
			},
			{
				Time:        lastRide.Time,
				Destination: lastRide.Origin,
			},
			{
				Time:        lastRide.Time,
				Destination: lastRide.Destination,
			},
		},
	}

	// Return the ride history of the passenger with the last ride repeated 3 times
	outputByte, err := json.Marshal(rideHistory)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %s", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(outputByte)
	return
}
