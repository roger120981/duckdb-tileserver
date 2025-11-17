package data

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// testEquals fails the test if exp is not equal to act.
func testEquals(tb testing.TB, exp, act interface{}, msg string) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: %s - expected: %#v; got: %#v\n", filepath.Base(file), line, msg, exp, act)
		tb.FailNow()
	}
}
