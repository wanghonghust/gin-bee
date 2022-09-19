package main

import "fmt"

type NumStr interface {
	string | Num
}

type Num interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func PrintNumStr[T NumStr](numStr T) {
	fmt.Printf("%T\n", numStr)
}

func main() {
	PrintNumStr(9)
}
