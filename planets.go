package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/soniakeys/meeus/v3/julian"
	"github.com/soniakeys/meeus/v3/moonposition"
	"github.com/soniakeys/meeus/v3/planetposition"
	"math"
	"time"
)

// Planet represents position and metadata of planet
type Planet struct {
	name string
	code int

	sB, cB, sL, cL float64 // sinus and cosinus of the planet
	r              float64 // heliocentric range in AU
}

// recalculatePositions recalculates positions of planets and moon
// and updates Prometheus metrics
func recalculatePositions(monitor *Monitor) error {
	// the calculation is based on the current time
	now := time.Now()
	jde := julian.TimeToJD(now)

	// report moon distance
	_, _, distance := moonposition.Position(jde)
	monitor.PlanceDistanceHistogram.With(prometheus.Labels{
		"from": "earth",
		"to":   "moon",
	}).Set(distance)

	// planets are sorted by the distance from the Sun
	planets := []*Planet{
		{name: "mercury", code: planetposition.Mercury},
		{name: "venus", code: planetposition.Venus},
		{name: "earth", code: planetposition.Earth},
		{name: "mars", code: planetposition.Mars},
		{name: "jupiter", code: planetposition.Jupiter},
		{name: "saturn", code: planetposition.Saturn},
		{name: "uranus", code: planetposition.Uranus},
		{name: "neptune", code: planetposition.Neptune},
	}

	// prepare all data needed for calculations
	// we want to calculate required values only once
	for _, p := range planets {
		planetData, err := planetposition.LoadPlanet(p.code)
		if err != nil {
			return fmt.Errorf("could not load planet data: %w", err)
		}
		L, B, R := planetData.Position(jde)
		sB, cB := B.Sincos()
		sL, cL := L.Sincos()

		p.sB = sB
		p.cB = cB
		p.sL = sL
		p.cL = cL
		p.r = R
	}

	// calculate and report distance between planets
	// do not calculate distance between the same planets twice
	// you can always swap from and to
	// planets are sorted by the distance from the Sun
	for k, p1 := range planets {
		for i := k + 1; i < len(planets); i++ {
			p2 := planets[i]

			// calculate distance between planets
			x := p1.r*p1.cB*p1.cL - p2.r*p2.cB*p2.cL
			y := p1.r*p1.cB*p1.sL - p2.r*p2.cB*p2.sL
			z := p1.r*p1.sB - p2.r*p2.sB
			delta := math.Sqrt(x*x + y*y + z*z)

			monitor.PlanceDistanceHistogram.With(prometheus.Labels{
				"from": p1.name,
				"to":   p2.name,
			}).Set(auToKm(delta))
		}
	}

	// calculate and report distance between planets and the Sun
	for _, p := range planets {
		monitor.PlanceDistanceHistogram.With(prometheus.Labels{
			"from": p.name,
			"to":   "sun",
		}).Set(auToKm(p.r))
	}

	return nil
}

// auToKm converts AU to km
// au is astronomical unit
func auToKm(au float64) float64 {
	return math.Round(au * 1.5 * math.Pow10(8))
}
