package evaluator

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"reflect"
	"time"
)

func Evaluate(expression string, timeout time.Duration) (float64, error) {
	time.Sleep(timeout)

	exp, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return 0, err
	}
	result, err := exp.Evaluate(nil)
	if err != nil {
		return 0, err
	}

	return getFloat(result, reflect.TypeOf(float64(0)))
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
