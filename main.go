package main

import "fmt"

func main() {
	var a, b int
	fmt.Scanln(&a, &b)
	result := add(a, b)
	fmt.Println(result)
}

func add(a int, b int) int {
	return a + b
}
