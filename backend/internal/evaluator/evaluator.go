package evaluator

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/k6mil6/distributed-calculator/backend/internal/response"
	"reflect"
	"time"
)

type Result struct {
	Id     int64   `json:"id"`
	Result float64 `json:"result"`
}

func Evaluate(response response.Response, heartbeat time.Duration, ch chan int64, workerId int64) (Result, error) {
	ticker := time.NewTicker(heartbeat)

	go func() {
		for range ticker.C {
			ch <- workerId
		}
	}()

	time.Sleep(response.Timeout)

	exp, err := govaluate.NewEvaluableExpression(response.Expression)
	if err != nil {
		return Result{}, err
	}
	result, err := exp.Evaluate(nil)
	if err != nil {
		return Result{}, err
	}

	resFloat, err := getFloat(result, reflect.TypeOf(float64(0)))
	if err != nil {
		return Result{}, err
	}

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
