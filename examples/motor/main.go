package main

import (
	"github.com/ricallinson/engine"
	"time"
)

func main() {

	// Start the robots processing engine.
	var robot = engine.Start(false)

	// Create a new instance of a Motor.
	var motor = robot.NewMotor(1, 2, false)

	// Spin the motor in a forwards direction and sleep for 5 seconds.
	motor.Forwards()
	time.Sleep(time.Second * 5)

	// Spin the motor in a backwards direction and sleep for 5 seconds.
	motor.Backwards()
	time.Sleep(time.Second * 5)

	// Stop the motor.
	motor.Stop()

	// At the end of the program it's good practice to stop the processing engine.
	robot.Stop()
}
