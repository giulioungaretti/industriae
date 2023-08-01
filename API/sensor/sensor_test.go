package sensor_test

import (
	"testing"
	"time"

	"atlan3d/api/sensor"
)

func TestSensor_Run(t *testing.T) {
	// create mock control and read channels
	// set the scan rate to 1 irco seconds
	// this is a cheat to predict the value of the sensor
	// as it will take one tick
	scanTime := 1
	setpoint := 10
	s := sensor.Create("test", scanTime)
	// set point
	s.ControlCh <- setpoint
	//
	expected := <-s.ReadCh
	if expected.SetPoint() != setpoint {
		t.Errorf("unexpected setpoint: %v", expected.SetPoint())
	}
	// wait for the sensor to read the value
	time.Sleep(100 * time.Millisecond)
	expected = <-s.ReadCh
	value := expected.Value()
	if value != setpoint {
		t.Errorf("unexpected value: %v", value)
	}
}
