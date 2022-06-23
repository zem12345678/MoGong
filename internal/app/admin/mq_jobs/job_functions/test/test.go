package test

import "fmt"

func Hello(req interface{}) (interface{}, error) {
	fmt.Println(req)
	return "hello test mq func", nil
}
