package evaluator

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/k6mil6/distributed-calculator/backend/internal/agent/response"
	"log/slog"
	"reflect"
	"time"
)

type Result struct {
	Id     int     `json:"id"`
	Result float64 `json:"result"`
}

func Evaluate(response response.Response, heartbeat time.Duration, ch chan int, workerId int, logger *slog.Logger) (Result, error) {
	ticker := time.NewTicker(heartbeat)
	defer ticker.Stop()

	logger.Info("evaluating expression", slog.Int("worker_id", workerId))

	go func() {
		for range ticker.C {
			ch <- workerId
		}
	}()

	time.Sleep(time.Duration(response.Timeout) * time.Second)

	logger.Info("expression timeout has gone", slog.Int("worker_id", workerId), slog.Any("expression", response.Subexpression))

	exp, err := govaluate.NewEvaluableExpression(response.Subexpression)

	if err != nil {
		return Result{}, err
	}
	result, err := exp.Evaluate(nil)

	logger.Info("expression evaluated", slog.Int("worker_id", workerId), slog.Any("result: %v", result))

	if err != nil {
		return Result{}, err
	}

	resFloat, err := getFloat(result, reflect.TypeOf(float64(0)))
	if err != nil {
		return Result{}, err
	}

	logger.Info("expression evaluated", slog.Int("worker_id", workerId), slog.Any("result: %v", resFloat))

	return Result{Id: response.Id, Result: resFloat}, nil
}

func getFloat(unk interface{}, floatType reflect.Type) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}
