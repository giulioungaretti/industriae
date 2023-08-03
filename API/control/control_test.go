package control_test

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"atlan3d/api/control"
	"atlan3d/api/sensor"
)

func TestLoop(t *testing.T) {
	// create a context with a timeout of 1 second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// create a data channel
	dataCh := make(chan sensor.SensorValues)

	// create two sensors with different set points
	sensor1 := sensor.Create("sensor1", 1, 100)
	sensor2 := sensor.Create("sensor2", 1, 100)
	sensor1SP := 50
	sensor2SP := 100
	sensor1.ControlCh <- sensor1SP
	sensor2.ControlCh <- sensor2SP

	// start the loop
	go control.Loop(ctx, 1, dataCh, sensor1, sensor2)

	// read data from the data channel
	var dataA, dataB sensor.SensorValues
	select {
	case dataA = <-dataCh:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for data from sensor")
	}
	select {
	case dataB = <-dataCh:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for data from sensor")
	}

	// check that the data is correct
	if dataA.Sensor() == "sensor1" {
		if dataA.SetPoint() != sensor1SP {
			t.Errorf("unexpected set point: %v", dataA.SetPoint())
		}
	} else {
		if dataA.SetPoint() != sensor2SP {
			t.Errorf("unexpected set point: %v", dataA.SetPoint())
		}
	}
	if dataB.Sensor() == "sensor1" {
		if dataB.SetPoint() != sensor1SP {
			t.Errorf("unexpected set point: %v", dataB.SetPoint())
		}
	} else {
		if dataB.SetPoint() != sensor2SP {
			t.Errorf("unexpected set point: %v", dataB.SetPoint())
		}
	}
}

func TestLoopTimeout(t *testing.T) {
	// create a context with a timeout of 1 second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// create a data channel
	dataCh := make(chan sensor.SensorValues)

	// create two sensors with different set points
	sensor1 := sensor.Create("sensor1", 1, 100)
	sensor2 := sensor.Create("sensor2", 1, 100)

	// start the loop
	go control.Loop(ctx, 10, dataCh, sensor1, sensor2)

	data := make([]sensor.SensorValues, 0)
outer:
	for {
		select {
		case dataA := <-dataCh:
			data = append(data, dataA)
		case <-ctx.Done():
			break outer
		default:
			continue
		}
	}
	expectedMinData := 100
	receivedLen := len(data)
	if receivedLen < expectedMinData {
		t.Errorf("unexpected number of data points: %v", receivedLen)
	}
}

func TestCsvLogger(t *testing.T) {
	// create a context with a timeout of 1 second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// create a data channel
	dataCh := make(chan sensor.SensorValues, 1)

	// create a temporary file for testing
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	sp := 1
	// send some data to the data channel
	data := sensor.Create("sensor1", 50, 100)
	data.ControlCh <- sp
	randomEmptyValue := <-data.ReadCh

	// read to chanenll
	dataCh <- randomEmptyValue
	// start the logger
	go control.CsvLogger(tmpfile.Name(), ctx, dataCh)

	// wait for the data to be processed
	time.Sleep(30 * time.Millisecond)

	dataCh <- randomEmptyValue
	// wait for the data to be processed
	time.Sleep(30 * time.Millisecond)

	// read the contents of the file
	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	// check that the data is correct
	if len(rows) < 2 {
		t.Fatalf("unexpected number of rows: %d", len(rows))
	}
	if rows[1][0] != "sensor1" {
		t.Errorf("unexpected sensor name: %s", rows[1][0])
	}
	if rows[1][1] != fmt.Sprintf("%d", randomEmptyValue.SetPointTs()) {
		t.Errorf("unexpected set point timestamp: %s", rows[1][1])
	}
	if rows[1][2] != fmt.Sprintf("%d", sp) {
		t.Errorf("unexpected set point: %s", rows[1][2])
	}
	if rows[1][3] != fmt.Sprintf("%d", sp) {
		t.Errorf("unexpected value: %s", rows[1][3])
	}
	if rows[1][4] != "null" {
		t.Errorf("unexpected error: %s", rows[1][4])
	}
}

func TestMultiplexWithPriority(t *testing.T) {
	// create a context with a timeout of 2 second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// create a data channel and two consumer channels
	dataCh := make(chan sensor.SensorValues, 1)
	consumer1 := make(chan sensor.SensorValues, 1)
	consumer2 := make(chan sensor.SensorValues, 1)

	// start the multiplexer
	control.MultiplexWithPriority(ctx, dataCh, consumer1, consumer2)

	// send some data to the data channel
	data := sensor.Create("sensor1", 50, 100)
	randomEmptyValue := <-data.ReadCh

	dataCh <- randomEmptyValue

	// wait for the data to be processed
	time.Sleep(60 * time.Microsecond)

	// check that the data was sent to the main channel and one of the consumer channels
	select {
	case <-consumer1:
		// data was sent to consumer1
	case <-consumer2:
		// data was sent to consumer2
	default:
		t.Error("data was not sent to any consumer")
	}
}
