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

type SolarSystem struct {
	planets map[int]*planetposition.V87Planet
	monitor *Monitor
}

func NewSolarSystem(monitor *Monitor) (*SolarSystem, error) {
	loaded := 0 // number of successfully loaded planets
	system := &SolarSystem{
		planets: map[int]*planetposition.V87Planet{
			planetposition.Mercury: loadPlanet(planetposition.Mercury, &loaded),
			planetposition.Venus:   loadPlanet(planetposition.Venus, &loaded),
			planetposition.Earth:   loadPlanet(planetposition.Earth, &loaded),
			planetposition.Mars:    loadPlanet(planetposition.Mars, &loaded),
			planetposition.Jupiter: loadPlanet(planetposition.Jupiter, &loaded),
			planetposition.Saturn:  loadPlanet(planetposition.Saturn, &loaded),
			planetposition.Uranus:  loadPlanet(planetposition.Uranus, &loaded),
			planetposition.Neptune: loadPlanet(planetposition.Neptune, &loaded),
		},
		monitor: monitor,
	}
	if loaded != len(system.planets) {
		return nil, fmt.Errorf("failed to load all planets data. You should run `make download` and set `VSOP87` environment variable")
	}
	return system, nil
}

// loadPlanet loads planet data from VSOP87B files
// loaded counter is a pointer, so we can increment it
func loadPlanet(planetCode int, loaded *int) *planetposition.V87Planet {
	planet, err := planetposition.LoadPlanet(planetCode)
	if err == nil {
		*loaded++
	}
	return planet
}

// recalculatePositions recalculates positions of planets and moon
// and updates Prometheus metrics
func (ss *SolarSystem) recalculatePositions() error {
	// Planet represents position and metadata of planet
	type Planet struct {
		name string
		code int

		sB, cB, sL, cL float64 // sinus and cosinus of the planet
		r              float64 // heliocentric range in AU
	}

	// the calculation is based on the current time
	now := time.Now()
	jde := julian.TimeToJD(now)

	// report moon distance
	_, _, distance := moonposition.Position(jde)
	ss.monitor.PlanceDistanceHistogram.With(prometheus.Labels{
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
	for _, p := range planets {
		planetData := ss.planets[p.code]
		L, B, R := planetData.Position(jde)
		p.sB, p.cB = B.Sincos()
		p.sL, p.cL = L.Sincos()
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

			ss.monitor.PlanceDistanceHistogram.With(prometheus.Labels{
				"from": p1.name,
				"to":   p2.name,
			}).Set(auToKm(delta))
		}
	}

	// calculate and report distance between planets and the Sun
	for _, p := range planets {
		ss.monitor.PlanceDistanceHistogram.With(prometheus.Labels{
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
