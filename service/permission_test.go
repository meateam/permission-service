package service_test

import (
	"testing"
	"math"
)

// Global variables

// Constants

func init() {
	
}

func Test(t *testing.T) {
	got := math.Abs(-1)
    if got != 1 {
        t.Errorf("Abs(-1) = %g; want 1", got)
    }
}
