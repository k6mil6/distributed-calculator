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
	id int
}

func New() *Worker {
	return &Worker{}
}

func (w *Worker) Start(url string, logger *slog.Logger, timeout, heartbeat time.Duration) {
	for {
		resp, err := http.Get(url + "/free_expressions")
		if err != nil {
			logger.Error("error sending POST request:", err)
			time.Sleep(timeout)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var mathResp response.Response

			logger.Info("expression received")

			err := json.NewDecoder(resp.Body).Decode(&mathResp)
			if err != nil {
				logger.Error("error decoding JSON response:", err)
				time.Sleep(timeout)
				continue
			}

			w.id = mathResp.Id

			logger.Info("expression decoded", slog.Int("worker_id", w.id), slog.Any("expression", mathResp.Subexpression))

			ch := make(chan int)

			go func() {
				err := sendHeartbeat(url, ch)
				if err != nil {
					logger.Error("error sending heartbeat:", err, "worker ID:", w.id)
				}
			}()

			res, err := evaluator.Evaluate(mathResp, heartbeat, ch, w.id, logger)
			if err != nil {
				logger.Error("error evaluating expression:", err)
				time.Sleep(timeout)
				continue
			}

			logger.Info("expression evaluated", slog.Int("worker_id", w.id))

			if err := sendResult(res, url); err != nil {
				logger.Error("error sending result:", err)
				time.Sleep(timeout)
				continue
			}

			logger.Info("result sent", slog.Int("worker_id", w.id))

		} else {
			logger.Error("non-OK response:", resp.StatusCode)
			time.Sleep(timeout)
			continue
		}

		resp.Body.Close()
		time.Sleep(timeout)
	}
}

func sendResult(result evaluator.Result, url string) error {
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

func sendHeartbeat(url string, ch <-chan int) error {
	for data := range ch {
		resp, err := http.Post(url+"/heartbeat", "application/json", bytes.NewBuffer(intToBytes(data)))
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

func intToBytes(i int) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}
