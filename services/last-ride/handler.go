package function

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
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

	// Add a sleep time for simulating database connection
	time.Sleep(time.Duration(rand.Intn(10)+5) * time.Millisecond)

	lastRide := Ride{
		PassengerID: *input.UserID,
	}
	switch *input.UserID {
	case 10:
		lastRide = Ride{
			Time: time.Now().Add(-10 * time.Minute),
			Origin: Point{
				Lat: 10.10,
				Lon: 40.40,
			},
			Destination: Point{
				Lat: 20.20,
				Lon: 30.30,
			},
		}
		break
	case 20:
		lastRide = Ride{
			Time: time.Now().Add(-20 * time.Minute),
			Origin: Point{
				Lat: 20.20,
				Lon: 30.30,
			},
			Destination: Point{
				Lat: 30.30,
				Lon: 20.20,
			},
		}
		break
	case 30:
		lastRide = Ride{
			Time: time.Now().Add(-30 * time.Minute),
			Origin: Point{
				Lat: 30.30,
				Lon: 20.20,
			},
			Destination: Point{
				Lat: 40.40,
				Lon: 10.10,
			},
		}
		break
	case 40:
		lastRide = Ride{
			Time: time.Now().Add(-40 * time.Minute),
			Origin: Point{
				Lat: 40.40,
				Lon: 10.10,
			},
			Destination: Point{
				Lat: 10.10,
				Lon: 40.40,
			},
		}
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("can not find the user"))
		return
	}

	outputByte, err := json.Marshal(lastRide)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %s", err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(outputByte)
}
