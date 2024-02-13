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
		numberData := map[string]int64{"worker_id": w.ID}
		jsonValue, _ := json.Marshal(numberData)

		resp, err := http.Post(url+"/freeExpressions", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			logger.Error("error sending POST request:", err)
			time.Sleep(timeout)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var mathResp response.Response

			logger.Info("expression received", slog.Int64("worker_id", w.ID))

			err := json.NewDecoder(resp.Body).Decode(&mathResp)
			if err != nil {
				logger.Error("error decoding JSON response:", err)
				time.Sleep(timeout)
				continue
			}

			logger.Info("expression decoded", slog.Int64("worker_id", w.ID), slog.Any("expression", mathResp.Subexpression))

			ch := make(chan int64)

			go func() {
				err := w.sendHeartbeat(url, ch)
				if err != nil {
					logger.Error("error sending heartbeat:", err, "worker ID:", w.ID)
				}
			}()

			res, err := evaluator.Evaluate(mathResp, heartbeat, ch, w.ID, logger)
			if err != nil {
				logger.Error("error evaluating expression:", err)
				time.Sleep(timeout)
				continue
			}

			logger.Info("expression evaluated", slog.Int64("worker_id", w.ID))

			if err := w.sendResult(res, url); err != nil {
				logger.Error("error sending result:", err)
				time.Sleep(timeout)
				continue
			}

			logger.Info("result sent", slog.Int64("worker_id", w.ID))

		} else {
			logger.Error("non-OK response:", resp.StatusCode)
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
