package testsuite

import (
	"fmt"
	"testing"
)

func ExampleIsDateTime() {
	fmt.Println(IsDateTime("2008-01-23T04:56:22Z"))
	// Output:
	// true
}

func TestIsBase64(t *testing.T) {
	for _, example := range []string{
		"",
		"Zg==",
		"Zm9v",
		"Zm9vYg==",
		"Zm9vYmE=",
		"Zm9vYmFy",
		"MY======", // base 32
		"MZXQ====", // base 32
	} {
		if !IsBase64(example) {
			t.Error(example)
		}
	}
}
