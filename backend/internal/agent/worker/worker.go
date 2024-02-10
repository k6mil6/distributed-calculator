package worker

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/k6mil6/distributed-calculator/backend/internal/agent/evaluator"
	"github.com/k6mil6/distributed-calculator/backend/internal/agent/response"
	"log/slog"
	"net/http"
	"time"
)

type Worker struct {
	ID int64
}

func New(id int64) *Worker {
	return &Worker{
		ID: id,
	}
}

func (w *Worker) Start(url string, logger *slog.Logger, timeout, heartbeat time.Duration) {
	for {
		resp, err := http.Get(url)
		if err != nil {
			logger.Error("Error sending GET request:", err)
			time.Sleep(timeout)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var mathResp response.Response
			err := json.NewDecoder(resp.Body).Decode(&mathResp)
			if err != nil {
				logger.Error("Error decoding JSON response:", err)
				time.Sleep(timeout)
				continue
			}

			ch := make(chan int64)

			go func() {
				err := w.sendHeartbeat(url, ch)
				if err != nil {
					logger.Error("Error sending heartbeat:", err, "Worker ID:", w.ID)
				}
			}()

			res, err := evaluator.Evaluate(mathResp, heartbeat, ch, w.ID)
			if err != nil {
				logger.Error("Error evaluating expression:", err)
				time.Sleep(timeout)
				continue
			}

			if err := w.sendResult(res, url); err != nil {
				logger.Error("Error sending result:", err)
				time.Sleep(timeout)
				continue
			}

		} else {
			logger.Error("Non-OK response:", resp.StatusCode)
			time.Sleep(timeout)
			continue
		}

		resp.Body.Close()
		time.Sleep(timeout)
	}
}

func (w *Worker) sendResult(result evaluator.Result, url string) error {
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := http.Post(url+"/result", "application/json", bytes.NewBuffer(jsonResult))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK response: %d", resp.StatusCode)
	}

	return nil
}

func (w *Worker) sendHeartbeat(url string, ch <-chan int64) error {
	for data := range ch {
		resp, err := http.Post(url+"/heartbeat", "application/json", bytes.NewBuffer(int64ToBytes(data)))
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("non-OK response: %d", resp.StatusCode)
		}
		resp.Body.Close()
	}

	return nil
}

func int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}
