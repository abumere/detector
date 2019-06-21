package travel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDistance(t *testing.T) {

	lat1, lon1 := 30.3773, -97.71
	lat2, lon2 := 34.0549, -118.2578

	distance := Distance(lat1, lon1, lat2, lon2)
	assert.Equal(t, 1226, int(distance * 0.000621371), "they should be equal")
}