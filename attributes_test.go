package testsuite

import (
	"fmt"
	"github.com/di-wu/regen"
	"testing"
)

func ExampleIsAttributeName() {
	fmt.Println(IsAttributeName("attrName"))
	// Output:
	// true
}

func TestIsAttributeName(t *testing.T) {
	g, err := regen.New(`[A-Za-z][\$\-_A-Za-z]*`)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 1000; i++ {
		if attrName := g.Generate(); !IsAttributeName(attrName) {
			t.Error(attrName)
		}
	}
}
