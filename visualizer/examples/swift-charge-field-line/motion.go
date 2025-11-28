package main

import (
	"math"

	"github.com/euphoricrhino/go-common/graphix"
)

type motion interface {
	pos(t float64, pos *graphix.Vec3)
	velAcc(t float64, vel *graphix.Vec3, acc *graphix.Vec3)
	frameToTime(f int) float64
}

type circularMotion struct {
	radius         float64
	omega          float64
	omegar         float64
	omega2r        float64
	initPhase      float64
	axis           int
	framesPerCycle int
}

func newCircularMotion(radius, omega, initPhase float64, axis, framesPerCycle int) *circularMotion {
	return &circularMotion{
		radius:         radius,
		omega:          omega,
		omegar:         omega * radius,
		omega2r:        omega * omega * radius,
		initPhase:      initPhase,
		axis:           axis,
		framesPerCycle: framesPerCycle,
	}
}

func (cm *circularMotion) pos(t float64, pos *graphix.Vec3) {
	arg := cm.omega*t + cm.initPhase
	pos[(cm.axis+1)%3] = cm.radius * math.Cos(arg)
	pos[(cm.axis+2)%3] = cm.radius * math.Sin(arg)
	pos[cm.axis%3] = 0.0
}

func (cm *circularMotion) velAcc(t float64, vel *graphix.Vec3, acc *graphix.Vec3) {
	arg := cm.omega*t + cm.initPhase
	cosArg, sinArg := math.Cos(arg), math.Sin(arg)
	vel[(cm.axis+1)%3] = -cm.omegar * sinArg
	vel[(cm.axis+2)%3] = cm.omegar * cosArg
	vel[cm.axis%3] = 0.0

	acc[(cm.axis+1)%3] = -cm.omega2r * cosArg
	acc[(cm.axis+2)%3] = -cm.omega2r * sinArg
	acc[cm.axis%3] = 0.0
}

func (cm *circularMotion) frameToTime(f int) float64 {
	return (2 * math.Pi / cm.omega) * float64(f) / float64(cm.framesPerCycle)
}

type harmonicMotion struct {
	amplitude float64
	omega     float64

	omegaa         float64
	omega2a        float64
	initPhase      float64
	axis           int
	framesPerCycle int
}

func newHarmonicMotion(
	amplitude, omega, initPhase float64,
	axis, framesPerCycle int,
) *harmonicMotion {
	return &harmonicMotion{
		amplitude:      amplitude,
		omega:          omega,
		omegaa:         omega * amplitude,
		omega2a:        omega * omega * amplitude,
		initPhase:      initPhase,
		axis:           axis,
		framesPerCycle: framesPerCycle,
	}
}

func (hm *harmonicMotion) pos(t float64, pos *graphix.Vec3) {
	arg := hm.omega*t + hm.initPhase
	pos[hm.axis%3] = hm.amplitude * math.Cos(arg)
	pos[(hm.axis+1)%3] = 0.0
	pos[(hm.axis+2)%3] = 0.0
}

func (hm *harmonicMotion) velAcc(t float64, vel *graphix.Vec3, acc *graphix.Vec3) {
	arg := hm.omega*t + hm.initPhase
	vel[hm.axis%3] = -hm.omegaa * math.Sin(arg)
	vel[(hm.axis+1)%3] = 0.0
	vel[(hm.axis+2)%3] = 0.0

	acc[hm.axis%3] = -hm.omega2a * math.Cos(arg)
	acc[(hm.axis+1)%3] = 0.0
	acc[(hm.axis+2)%3] = 0.0
}

func (hm *harmonicMotion) frameToTime(f int) float64 {
	return (2 * math.Pi / hm.omega) * float64(f) / float64(hm.framesPerCycle)
}
