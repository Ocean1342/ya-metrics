package handlers

import (
	"fmt"
)

func ExampleUpdateRequestPrepare() {
	path1 := "localhost:8080/update/counter/PollCount/1"
	ur1, _ := UpdateRequestPrepare(path1)
	fmt.Println(*ur1)

	path2 := "localhost:8080/update/gauge/Alloc/1.23456"
	ur2, _ := UpdateRequestPrepare(path2)
	fmt.Println(*ur2)

	// Output:
	// {counter PollCount 1}
	// {gauge Alloc 1.23456}
}
