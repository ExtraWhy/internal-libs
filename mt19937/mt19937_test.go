package mt19937

import (
	"fmt"
	"testing"
)

func TestInitEX(t *testing.T) {
	InitEX()
	for i := 0; i < 1000; i++ {
		fmt.Println(Genrand64_int64() % 34)
	}
}
