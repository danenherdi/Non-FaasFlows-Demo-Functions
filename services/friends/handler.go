package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Input is the argument of your flow function
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

	// Request to the user info
	reqBody, err := json.Marshal(map[string]interface{}{
		"user_id": *input.UserID,
	})
	req, err := http.NewRequest(
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
	client := &http.Client{}
	resp, err := client.Do(req)
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
	data, err := io.ReadAll(resp.Body)
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

	// Create friends info
	friendsInfo := FriendsInfo{
		NumberOfFriends: userInfo.PhoneNumber,
	}

	outputByte, err := json.Marshal(friendsInfo)
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
