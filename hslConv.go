package main

import (
	"fmt"
	"math"
)

type RGB struct {
	r float64
	g float64
	b float64
}

func HUEtoRGB(p float64, q float64, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}

	if t < 1./6. {
		return p + (q-p)*6*t
	} else if t < 1./2. {
		return q
	} else if t < 2./3. {
		return p + (q-p)*(2./3.-t)*6
	} else {
		return p
	}
}

func RGBtoHEX(val float64) string {
	hex := fmt.Sprintf("%x", int(math.Round(val*255)))
	if len(hex) == 1 {
		return "0" + hex
	}
	return hex
}

func HSLtoHEX(h float64, s float64, l float64) string {
	h /= 360
	s /= 100
	l /= 100

	var rgb RGB

	if s == 0 {
		rgb.r = l
		rgb.b = l
		rgb.g = l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		rgb.r = HUEtoRGB(p, q, h+1./3.)
		rgb.g = HUEtoRGB(p, q, h)
		rgb.b = HUEtoRGB(p, q, h-1./3.)
	}

	return "#" + RGBtoHEX(rgb.r) + RGBtoHEX(rgb.g) + RGBtoHEX(rgb.b)
}
