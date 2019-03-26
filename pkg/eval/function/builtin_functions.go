//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package function

import (
	"math"
	"reflect"

	"github.com/oracle/speedle/pkg/errors"
)

// Add all built-in functions in this file

func Sqrt(args ...interface{}) (interface{}, error) {
	err := errors.New(errors.BuiltInFuncError, "Usage: Sqrt(x)")
	if len(args) != 1 {
		return nil, err
	}
	if reflect.TypeOf(args[0]).Kind() != reflect.Float64 {
		return nil, err
	}
	x := args[0].(float64)
	return math.Sqrt(x), nil
}

func Max(args ...interface{}) (interface{}, error) {
	err := errors.New(errors.BuiltInFuncError, "Usage: Max(x1, x2, ...), xi must be numeric")
	if len(args) < 1 {
		return nil, err
	}
	var max = args[0].(float64)
	for i := range args[1:] {
		if reflect.TypeOf(args[i]).Kind() != reflect.Float64 {
			return nil, err
		}
		max = math.Max(max, args[i].(float64))
	}
	return max, nil
}

func Min(args ...interface{}) (interface{}, error) {
	err := errors.New(errors.BuiltInFuncError, "Usage: Min(x1, x2, ...), xi must be numeric")
	if len(args) < 1 {
		return nil, err
	}
	var min = args[0].(float64)
	for i := range args[1:] {
		if reflect.TypeOf(args[i]).Kind() != reflect.Float64 {
			return nil, err
		}
		min = math.Min(min, args[i].(float64))
	}
	return min, nil
}

func Sum(args ...interface{}) (interface{}, error) {
	err := errors.New(errors.BuiltInFuncError, "Usage: Sum(x1, x2, ...), xi must be numeric")
	var sum float64 = 0
	for i := range args {
		if reflect.TypeOf(args[i]).Kind() != reflect.Float64 {
			return nil, err
		}
		sum += args[i].(float64)
	}
	return sum, nil
}

func Avg(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return float64(0), nil
	}
	sum, err := Sum(args)
	if err != nil {
		return nil, err
	}
	return sum.(float64) / float64(len(args)), nil
}

// IsSubSet(S1, S2) means "is S1 a subset of S2"
func IsSubSet(args ...interface{}) (interface{}, error) {
	err := errors.New(errors.BuiltInFuncError, "Usage: IsSubSet(S1, S2) - S1 and S2 are both slice, and test if S1 is a subset of S2")
	n := len(args)
	if n < 2 {
		return nil, err
	}
	s1, s2 := args[0], args[n-1]
	if n >= 2 {
		buf := make([]interface{}, n-1)
		copy(buf, args)
		s1 = buf
	}
	if reflect.TypeOf(s1).Kind() != reflect.Slice || reflect.TypeOf(s2).Kind() != reflect.Slice {
		return nil, err
	}
	v1 := reflect.ValueOf(s1)
	v2 := reflect.ValueOf(s2)
	n1 := v1.Len()
	n2 := v2.Len()
	if n1 == 0 || n2 == 0 || n1 > n2 {
		return false, nil
	}
outer:
	for i := 0; i < n1; i++ {
		for j := 0; j < n2; j++ {
			if v1.Index(i).Interface() == v2.Index(j).Interface() {
				continue outer
			}
		}
		return false, nil
	}
	return true, nil
}
