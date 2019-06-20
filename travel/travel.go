package travel

import (
	"fmt"
	"math"
	"time"
)
////////////////////////////
////////////////////////////
// These functions were referenced from: https://gist.github.com/cdipaolo/d3f8db3848278b49db68
////////////////////////////
////////////////////////////
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func Speed(distance float64, startT, endT int64) int {
	distInMiles := distance *  0.00062137
	startTime := time.Unix(startT, 0)
	endTime := time.Unix(endT, 0)
	speed := distInMiles/math.Abs(endTime.Sub(startTime).Hours())

	//fmt.Println("Start Time: ", startTime)
	//fmt.Println("End Time: ", endTime)
	fmt.Println("Distance in Miles: ", distInMiles)
	fmt.Println("Time Difference: ", math.Abs(endTime.Sub(startTime).Hours()) )
	fmt.Println("Speed: ", int(speed), " miles per hour")

	//newPreceedingTime := time.Unix(1514764800,0).Add(time.Hour * -24).Add(time.Minute * -18).Add(time.Second * -41)
	//fmt.Println("New Time: ", newPreceedingTime.Unix())
	//fmt.Println("Difference between the two hours: ", endTime.Sub(startTime).Hours())
	return int(speed)

}