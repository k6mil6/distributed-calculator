package worker

import (
	"fmt"
	"github.com/k6mil6/distributed-calculator/backend/internal/evaluator"
	"io/ioutil"
	"net/http"
	"time"
)

type Worker struct {
	ID int64
}

func NewWorker(id int64) *Worker {
	return &Worker{
		ID: id,
	}
}

func (w *Worker) Start() {
	for {
		// Send GET request to the HTTP server
		resp, err := http.Get("http://example.com/endpoint")
		if err != nil {
			// Handle error
			fmt.Println("Error:", err)
			time.Sleep(1 * time.Second) // Wait and send GET request again
			continue
		}
		defer resp.Body.Close()

		// Process the response based on the content
		if resp.StatusCode == http.StatusOK {
			// Read the response body
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// Handle error
				fmt.Println("Error reading response body:", err)
				continue
			}
			// Process the response body (evaluate math expression)
			evaluator.Evaluate(string(body), time.Second)
		} else {
			// Handle non-OK response
			fmt.Println("Non-OK response:", resp.Status)
		}

		// Wait and send GET request again
		time.Sleep(1 * time.Second)
	}
}
