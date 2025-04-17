package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

	if input.UserID == nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("user id is empty"))
		return
	}

	if input.Origin == nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("origin is empty"))
		return
	}

	// Request to last ride
	reqBody, err := json.Marshal(map[string]interface{}{
		"user_id": *input.UserID,
	})
	req, err := http.NewRequest(
		"POST",
		"http://gateway.openfaas:8080/function/last-ride-nonflow",
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

	// Request to the user info
	// Request to last ride
	reqBody, err = json.Marshal(map[string]interface{}{
		"user_id": *input.UserID,
	})
	req, err = http.NewRequest(
		"POST",
		"http://gateway.openfaas:8080/function/user-info-nonflow",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		fmt.Printf("error in creating new request of function user-info-nonflow: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// Do the request
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("error in doing the request of function user-info-nonflow: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("failed request of user-info-nonflow")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("failed request of user-info-nonflow"))
		return
	}
	// Parse the response
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	var userInfo UserInfo
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("%+v\n", userInfo)

	var recommendation Recommendation
	// Repeat recommendation
	if input.Origin.Lat == lastRide.Origin.Lat && input.Origin.Lon == lastRide.Origin.Lon {
		recommendation = Recommendation{
			Type:           RecommendationRepeat,
			Recommendation: &lastRide.Destination,
			BannerText:     fmt.Sprintf("Dear %s, Here is your repeat recommendation.", userInfo.FirstName),
		}

		outputByte, err := json.Marshal(recommendation)
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

	// Reverse recommendation
	if input.Origin.Lat == lastRide.Destination.Lat && input.Origin.Lon == lastRide.Destination.Lon {
		recommendation = Recommendation{
			Type:           RecommendationReverse,
			Recommendation: &lastRide.Origin,
			BannerText:     fmt.Sprintf("Dear %s, Here is your reverse recommendation.", userInfo.FirstName),
		}

		outputByte, err := json.Marshal(recommendation)
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

	// No recommendation
	recommendation = Recommendation{
		Type:       RecommendationNothing,
		BannerText: fmt.Sprintf("Dear %s, There is no recommendation.", userInfo.FirstName),
	}

	outputByte, err := json.Marshal(recommendation)
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
