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
	setPointLimit := 10
	s := sensor.Create("test", scanTime, setPointLimit)
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

func TestSensor_Error(t *testing.T) {
	// create mock control and read channels
	// set the scan rate to 1 irco seconds
	// this is a cheat to predict the value of the sensor
	// as it will take one tick
	scanTime := 1
	setpoint := 10
	setPointLimit := 0
	s := sensor.Create("test", scanTime, setPointLimit)
	// set setpoint to much higher value
	s.ControlCh <- setpoint
	expected := <-s.ReadCh
	if expected.SetPoint() == setpoint {
		t.Errorf("unexpected setpoint: %v should have not ben set over the limit", expected.SetPoint())
	}
	if expected.Error() == nil {
		t.Errorf("Expected error")
	}
	// wait for the sensor to read the value
	time.Sleep(100 * time.Millisecond)
	expected = <-s.ReadCh
	value := expected.Value()
	if value == setpoint {
		t.Errorf("unexpected value: %v", value)
	}
}
