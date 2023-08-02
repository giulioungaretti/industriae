package sensor

import (
	"math"
	"strconv"
	"time"
)

type SensorError string

// Error returns the error message
func (e SensorError) Error() string {
	return string(e)
}

// ErrSetPointTooHigh is returned when the setpoint is higher than the limit
var ErrSetPointTooHigh = SensorError("setpoint too high")

type SensorValues struct {
	sensor     string
	setPoint   int
	limit      int
	setPointTs int64
	value      int
	ts         int64
	error      SensorError
}

// Timestamp is the timestamp of the last sensor read
func (s SensorValues) Ts() int64 {
	return s.ts
}

// Sensor is the name of the sensor
func (s SensorValues) Sensor() string {
	return s.sensor
}

// SetPoint is the current setpoint of the sensor
func (s SensorValues) SetPoint() int {
	return s.setPoint
}

// SetPointTs is the timestamp of the last setpoint
func (s SensorValues) SetPointTs() int64 {
	return s.setPointTs
}

// Value is the current sensor value
func (s SensorValues) Value() int {
	return s.value
}

// Error is the last error encountered by the sensor
func (s SensorValues) Error() error {
	return s.error
}

// Limit is the maximum value of the sensor
func (s SensorValues) Limit() int {
	return s.limit
}

// String returns the sensor values as a string array
func (s SensorValues) String() []string {
	error := "null"
	// if s.erorr != nil {
	// error = s.erorr.Error()
	// }
	return []string{
		s.sensor,
		strconv.FormatInt(s.setPointTs, 10),
		strconv.Itoa(s.setPoint),
		strconv.Itoa(s.value),
		error,
	}
}

type Sensor struct {
	SensorValues
	// readCh is used to read the sensor values
	ReadCh chan SensorValues
	// controlCh is used to set the setpoint of the sensor
	ControlCh chan int
	// scanRate is the rate at which the sensor is read in milliseconds
	scanRate int
}

func (s *Sensor) run() {
	for {
		select {
		case spm := <-s.ControlCh:
			if spm > s.limit {
				s.error = ErrSetPointTooHigh
			}
			s.SensorValues.setPointTs = time.Now().UnixNano()
			s.SensorValues.setPoint = spm
		case <-time.After(time.Duration(s.scanRate) * time.Millisecond):
			//	simulate a lag in the setpoint/set of the sensor
			// the logic is silly beacuse it is just a simulation
			if s.SensorValues.value <= s.SensorValues.setPoint {
				diff := math.Abs(float64(s.SensorValues.setPoint) - float64(s.SensorValues.value))
				sum := math.Floor(diff / 3)
				s.SensorValues.value += int(sum)
				if sum < 1 {
					s.SensorValues.value = s.SensorValues.setPoint
				}
			}
			if s.SensorValues.value >= s.SensorValues.setPoint {
				diff := math.Abs(float64(s.SensorValues.setPoint) - float64(s.SensorValues.value))
				sum := math.Floor(diff / 3)
				s.SensorValues.value -= int(sum)
				if sum < 1 {
					s.SensorValues.value = s.SensorValues.setPoint
				}
			}
			s.SensorValues.ts = time.Now().UnixNano()
			select {
			// sensor is not blocked, send the data
			case s.ReadCh <- s.SensorValues:
			//  sensor is not read, drop the data
			case <-s.ReadCh:
				// log.Trace().Msgf("dropping stale data from sensor %s", s.SensorValues.sensor)
			default:
				// continue
			}
		}
	}
}

// Create a new sensor with the given name
// starts the sensor loop in a goroutine with the given scan rate in microseconds
func Create(name string, scanRate int, limit int) Sensor {
	s := Sensor{
		// TODO: set the initial value to the current value of the sensor setting to 0 might be dangerous default
		SensorValues: SensorValues{sensor: name, setPoint: 0, value: 0, ts: time.Now().UnixNano(), limit: limit},
		// always buffer for two values, so that the sensor can be read while the control loop is running
		ReadCh: make(chan SensorValues, 2),
		// always buffer for one incoming command
		ControlCh: make(chan int, 1),
		scanRate:  scanRate,
	}
	go s.run()
	return s
}
