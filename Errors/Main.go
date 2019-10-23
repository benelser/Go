package main

import (
	. "fmt"
)

func main(){
	// Naive programmers fail to implement errors especially around IO operations.

	// If failure has only one possible cause the result is boolean usually called ok

	marriedCouple := make(map[int]string)
	marriedCouple[1] = "Amanda"
	marriedCouple[2] = "Benjamin"
	for key, value := range marriedCouple {
		Printf("%d %s\r\n", key, value)
	}

	// Check if a key exists or not
	value, ok := marriedCouple[3]
	if ok{
		Printf("Value of Key 3 is : %s\r\n", value)
	}

	value2, ok := marriedCouple[2]
	if ok{
		Printf("value of Key 2 is : %s\r\n", value2)
	}

	// This is the idiomatic way of handling errors in Go
	// if err != nil {
	// 	return nil, err
	// }
	
}