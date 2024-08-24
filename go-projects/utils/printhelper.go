package utils

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
)

func PrintQuery(query string, args ...interface{}) {
	callerFuncName, line := callerFunctionName()
	info := formatQuery(query, args...)
	log.Printf("[DB LOG]\n\tFunc: %s\n\tLine:%d\n\tQuery: %s\n", callerFuncName, line, info)
}

func callerFunctionName() (string, int) {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name(), line
}

func formatQuery(query string, args ...interface{}) string {
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		argStr := formatArg(arg)
		query = strings.Replace(query, placeholder, argStr, 1)
	}
	return query
}

func formatArg(arg interface{}) string {
	// Desreferenciar ponteiros
	value := reflect.ValueOf(arg)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return "NULL"
		}
		arg = value.Elem().Interface()
	}

	switch v := arg.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case fmt.Stringer:
		return fmt.Sprintf("'%s'", v.String())
	default:
		return fmt.Sprintf("%v", v)
	}
}
