//
// Copyright 2017, Yahoo Inc.
// Copyrights licensed under the New BSD License.
// See the accompanying LICENSE file for terms.
//

package engine

import (
	"github.com/ricallinson/engine/gpio"
	"log"
	"time"
)

type RangeSensor struct {
	pinTrigger *PinOutput
	pinEcho    *PinInput
}

// Returns a new instance of RangeSensor.
// The value of `pin` must be in the range of 1-25 mapping to the Raspberry Pi GPIO pins.
// Controls a HC-SR04 Ultrasonic Range Sensor.
func NewRangeSensor(pinTrigger *gpio.GpioPin, pinEcho *gpio.GpioPin) *RangeSensor {
	this := &RangeSensor{
		pinTrigger: NewPinOutput(pinTrigger),
		pinEcho:    NewPinInput(pinEcho),
	}
	// Wait for the sensor to settle.
	time.Sleep(2 * time.Second)
	return this
}

// Returns the pins that this instance is controlled by.
func (this *RangeSensor) Pins() (uint8, uint8) {
	return this.pinTrigger.PinOut(), this.pinEcho.PinOut()
}

// Returns a measurement between 1cm and 400cm.
func (this *RangeSensor) Get() float32 {
	// The HC-SR04 sensor requires a short 10uS pulse to trigger the module, which will
	// cause the sensor to start the ranging program (8 ultrasound bursts at 40 kHz) in
	// order to obtain an echo response. So, to create our trigger pulse, we set out
	// trigger pin high for 10uS then set it low again. Send a short rpio.Low pulse just to make sure.
	this.pinTrigger.Low()
	time.Sleep(5 * time.Microsecond)
	this.pinTrigger.High()
	time.Sleep(20 * time.Microsecond)
	this.pinTrigger.Low()
	// Measure the distance.
	distance := this.takeMeasurement()
	// If the measurement is out of range return 0.
	if distance < 1 {
		return -1
	}
	if distance > 400 {
		return 400
	}
	// Log the final measurement and return it.
	this.log(distance)
	return distance
}

// Returns the distance measured in cm units.
func (this *RangeSensor) takeMeasurement() float32 {
	// It should be about 200ms before the sensor has data so use this time to setup the measurement variables.
	var pulseStart = time.Now()
	var pulseDuration time.Duration
	// Our first step is to record the last rpio.Low timestamp for pinEcho (pulseStart)
	// e.g. just before the return signal is received and pinEcho goes rpio.High.
	for this.pinEcho.Read() == gpio.Low {
		// Allow somehting else to happen on a single core Raspberry Pi.
		time.Sleep(time.Nanosecond)
		// If more than 38ms was spent here the measurement failed.
		if time.Since(pulseStart).Seconds() > 1 {
			return -1
		}
	}
	// We are a go so recored this moment as the start.
	pulseStart = time.Now()
	// Once a signal is received, the value changes from rpio.Low (0) to rpio.High (1), and the
	// signal will remain rpio.High for the duration of the pinEcho pulse. We therefore also need
	// the last rpio.High timestamp for pinEcho to give us a duration.
	// We wait for that moment to come like a high school boy on a first date.
	for this.pinEcho.Read() == gpio.High {
		// Allow somehting else to happen on a single core Raspberry Pi.
		time.Sleep(time.Nanosecond)
		// If more than 38ms was spent here the measurement failed.
		if time.Since(pulseStart).Seconds() > 1 {
			return -1
		}
	}
	// Get the duration of the time taken for sound to travel to the obstacle and back.
	pulseDuration = time.Since(pulseStart)
	// The formula for distance measured on the HC-SR04 sensor is cm = uS / 58.
	return float32(pulseDuration.Nanoseconds() / 58000)
}

// Logs state of the assigned pin.
func (this *RangeSensor) log(cm float32) {
	log.Print("RangeSensor on pin ", this.pinTrigger, " and ", this.pinEcho, " measured a distance of ", cm, "cm")
}
