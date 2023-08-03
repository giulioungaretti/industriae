package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"atlan3d/api/control"
	"atlan3d/api/sensor"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins
		return true
	},
}

type SensorValuesMsg struct {
	Sensor     string
	SetPoint   int
	SetPointTs int64
	Value      int
	Ts         int64
	Error      error
	TsSent     int64
}

func main() {
	// set up logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// http server adress and port and context
	serverAddress := ":8080"
	// router httpserver
	r := mux.NewRouter()
	backgroundCtx := context.Background()

	// set up sensors
	pressureSensor := sensor.Create("Pressure1", 1, 1000)
	temperatureSensor := sensor.Create("Temperature1", 1, 100)
	sensors := []sensor.Sensor{pressureSensor, temperatureSensor}
	dataStreamCh := make(chan sensor.SensorValues)

	// status of the control loop
	status := "stopped"

	// channels to multiplex data to
	logChannel := make(chan sensor.SensorValues)
	wsChannel := make(chan sensor.SensorValues, 100)

	// a slice of channels to multiplex the data from the sensors to the main channel and other consumsers
	multiplexDataChannels := []chan sensor.SensorValues{logChannel, wsChannel}
	control.MultiplexWithPriority(backgroundCtx, dataStreamCh, multiplexDataChannels...)

	// create child context
	var controlCtx context.Context
	var controlCtxCancelFunction context.CancelFunc
	r.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		if status == "stopped" {
			controlCtx, controlCtxCancelFunction = context.WithCancel(backgroundCtx)
			// start the control loop
			control.Loop(controlCtx, 50, dataStreamCh, sensors...)
			// start the csv logger
			go control.CsvLogger("data.csv", backgroundCtx, logChannel)

			status = "running"

			// response
			json.NewEncoder(w).Encode(map[string]string{"status": "started"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "already started"})
	})

	r.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if status == "stopped" {
			json.NewEncoder(w).Encode(map[string]string{"status": "already stopped"})
			return
		}
		controlCtxCancelFunction()
		status = "stopped"
		json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
	})

	r.HandleFunc("/temperature", func(w http.ResponseWriter, r *http.Request) {
		if status == "stopped" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "control loop is not running, please start it first"})
		}
		// set value
		tempStr := r.URL.Query().Get("value")
		if tempStr != "" {
			v, err := strconv.Atoi(tempStr)
			if err == nil {
				select {
				case temperatureSensor.ControlCh <- int(v):
					err := <-temperatureSensor.ErrorCh
					if err != nil {
						json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
						return
					}
					json.NewEncoder(w).Encode(map[string]int{"ok": temperatureSensor.Value()})
					return
				default:
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "control loop is not accepting commands, please try again later"})
					return
				}
			}
		}
	})

	r.HandleFunc("/pressure", func(w http.ResponseWriter, r *http.Request) {
		if status == "stopped" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "control loop is not running, please start it first"})
		}
		pressureStr := r.URL.Query().Get("value")
		if pressureStr != "" {
			v, err := strconv.Atoi(pressureStr)
			if err == nil {
				select {
				case pressureSensor.ControlCh <- int(v):
					json.NewEncoder(w).Encode(map[string]float64{"pressure": rand.Float64()})
					err := <-pressureSensor.ErrorCh
					if err != nil {
						json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
						return
					}
					json.NewEncoder(w).Encode(map[string]int{"ok": pressureSensor.Value()})
					return
				default:
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "control loop is not accepting commands, please try again later"})
					return
				}
			}
		}
	})

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if status == "stopped" {
			w.WriteHeader(http.StatusBadRequest)
			if status == "stopped" {
				json.NewEncoder(w).Encode(map[string]string{"message": "control loop is not running, please start it first"})
			}
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Panic().Msgf("upgrade socket failed: %s", err)
			return
		}

		go func(conn *websocket.Conn) {
			for {
				select {
				case message := <-wsChannel:
					data := SensorValuesMsg{
						Sensor:     message.Sensor(),
						SetPoint:   message.SetPoint(),
						SetPointTs: message.SetPointTs(),
						Value:      message.Value(),
						Ts:         message.Ts(),
						Error:      message.Error(),
						TsSent:     time.Now().UnixNano(),
					}
					conn.WriteJSON(data)
				case <-controlCtx.Done():
					return
				default:
					log.Trace().Msgf("no message for ws client ")
					time.Sleep(10 * time.Millisecond)
				}
			}
		}(conn)
	})

	// Create a CORS handler with default options
	// TODO: this is very broad, only for dev purposes, should be more restrictive in prod
	corsHandler := cors.Default().Handler(r)

	log.Info().Msgf("Starting API server... %s", serverAddress)
	http.ListenAndServe(serverAddress, corsHandler)
}
