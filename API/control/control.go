package control

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	"atlan3d/api/sensor"

	"github.com/rs/zerolog/log"
)

// multiplexWithPriority takes a channel of any and multiplexes it to a main channel and a set of consumer channels.
// The main channel is used to process the message with the mainFunc, which is guaranteed to run, block, or error.
// messages to the consumer channels are dropped if the channel is blocked, or buffer is full.
func MultiplexWithPriority(
	ctx context.Context,
	input chan sensor.SensorValues,
	consumers ...chan sensor.SensorValues,
) {
	// We'll use a sync.Once to make sure we don't start a bunch of these.
	var once sync.Once
	main, rest := consumers[0], consumers[1:]
	once.Do(
		func() {
			go func() {
				// Every time a message comes over the channel...
				for v := range input {
					// Send it to the main channel
					main <- v
					// Loop over the  rest of consumers...
					for _, cons := range rest {
						// Send each one the message but drop them if the channels block
						select {
						case cons <- v:
							log.Trace().Msgf("Sending message %v", v)
							// we can send the message
						case <-ctx.Done():
							// we done, return
						default:
							// we can't send the message, since it's not the main channel we don't care
							// it's OK to drop the message
							log.Trace().Msgf("dropping message %v", v)
						}
					}
				}
			}()
		})
}

// csvLogger writes the data to a csv file with the following format
// ts, sensor, op, value, error
// NOTE: writing to one file is the fastest and safest way to log data
// NOTE: use e.g. duckdb in memory to read as strucctued data later
func CsvLogger(
	output string,
	ctx context.Context,
	dataCh chan sensor.SensorValues,
) {
	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// open a writer data to the CSV file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	header := []string{"sensor", "sp-ts", "sp", "value", "error"}
	writer.Write(header)
	writer.Flush()
	err = writer.Error()
	if err != nil {
		log.Error().Msgf("error flushing to csv file: %s", err.Error())
	}

	for {
		select {
		case data := <-dataCh:
			// Write the sensor data as a new row
			row := data.String()
			err := writer.Write(row)
			if err != nil {
				log.Error().Msgf("error writing to csv file: %s", err)
			}
		case <-ctx.Done():
			fmt.Println("Received stop signal. Stopping.")
			writer.Flush()
			err = writer.Error()
			if err != nil {
				log.Error().Msgf("error flushing to csv file: %s", err.Error())
			}
			return
		case <-time.After(10 * time.Millisecond):
			writer.Flush()
			err = writer.Error()
			if err != nil {
				log.Error().Msgf("error flushing to csv file: %s", err.Error())
			}
		}
	}
}

// Looop starts the loop over the sensors and reads data from them, and sends it to the data channel
func Loop(
	ctx context.Context,
	// scan rate in millisceonds
	scanFreq int,
	dataCh chan sensor.SensorValues,
	sensors ...sensor.Sensor,
) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msgf("Received stop signal. Stopping.")
				return
			case <-time.After(time.Duration(scanFreq) * time.Millisecond):
				// Loop over the sensors...
				for _, sensor := range sensors {
					select {
					case data := <-sensor.ReadCh:
						log.Debug().Msgf("Received data from sensor %s", sensor.SensorValues.Sensor())
						dataCh <- data
					case <-ctx.Done():
						log.Info().Msgf("Received stop signal. Stopping.")
					default:
						log.Trace().Msgf("no data from sensor %s in the last %d",
							sensor.SensorValues.Sensor(), scanFreq)
					}
				}
			}
		}
	}()
}
