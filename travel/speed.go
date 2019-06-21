package travel

import (
	"math"
	"time"
)

// Used to calculate the speed in mph of traveling a certain distance in a certain time
func Speed(distance float64, startT, endT int64) int {
	distInMiles := distance *  0.00062137
	startTime := time.Unix(startT, 0)
	endTime := time.Unix(endT, 0)
	speed := distInMiles/math.Abs(endTime.Sub(startTime).Hours())

	//fmt.Println("Distance in Miles: ", distInMiles)
	//fmt.Println("Time Difference: ", math.Abs(endTime.Sub(startTime).Hours()) )
	//fmt.Println("Speed: ", int(speed), " miles per hour")

	return int(speed)
}