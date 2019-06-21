package travel

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSpeed(t *testing.T) {

	startTime := time.Unix(1514851200, 0)
	endTime := startTime.Add(time.Hour * 10)
	//1000 miles
	distance := 1609347.08789

	speed := Speed(distance, startTime.Unix(), endTime.Unix())

	assert.Equal(t, 100, speed, "they should be equal")
}