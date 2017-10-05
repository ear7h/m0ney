package data

import (
	"testing"
	"fmt"
)

func TestGetSet(t *testing.T) {
	s := GetSet(1)
	fmt.Println(s)
}