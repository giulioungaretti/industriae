# industriae

A Industrial Control System prototype monorepo.

## Principles

Quality. All software must delight and work.
Safety.  All software must be safe.
Clarity. All software must have purpose and create value.

Based on this principles, a few choices were made:

*   No c/c++ unless documented exception or required (rtos/no llvm), if closer to metal is needed consider memory safe alternatives. E.g. rust.

*   The prototype will be scarce in features, but of high quality and clear purpose.

## Getting started
This monorepo contains the following projects:
- API - a go project/module that contains both the REST api and a simulated embedded system taht mimics the rtos.
- GUI - a react/typescript app that provides a simple UI to control the system and display the data.

It is **required** to install docker to run both compoments, and the vscode extension [dev container](https://code.visualstudio.com/docs/devcontainers/create-dev-container) !

Inside the two folders there are instructions on how to run the projects.  
It is **required** to run the API first, and then the GUI!

## Architecture


### The system

There is a oven controlled by an simulated embedded rtos, which exposes the api described below.

A control loop can be started and stopped. The control loop is configured to to read from a list of sensors (temperature and pressure) and multiplex back the data to a series of streams. The first stream is loseless, and it's meant to log all the data. The other streams drop data if there are no listeners or if they can't keep up. The lossy streams can be used to stream the data back to the ui for example.

The control loop can be configured to run at a certain frequency, and the sensors can be configured to run at a certain frequency.

The sensors are simulated, and they are configured to generate very simple dampened response to the setpoint they recieve.

Every action is non blocking and async, and the system is designed to be fault tolerant.

### User Interface

A simple UI allows users to start and stop the equipment, and set the setpoint of the temperaure and pressure sensors. The UI also displays the real-time data from the sensors.

![GUI](./GUI.png)

### Data Logging

Data logging is simple. The system should log all the data from the sensors. The data is stored in a simple flat file as csv data. This data can be used for further analysis and nothing is faster than an append only write to disk.
Log rotation is not implmented in this symple example as it relies on the unerlying OS.
Furhter analyics can be done in ohter systems by reading the LOG file(s) eg using duckdb for parsing and querying or shipping to other cloud systems.

### Safety Protocols

Basic saftey protcols are implmented as simple bound checks in the embedded system.
Rarely should we trust higher level systems to deal with hardware/safety limits.

### API Design

The API is designed to be simple and easy to use. The API is designed to be async and non blocking and based on REST principles, for live data a websocked connection is used. GOLANG is used becaue it's simple and fast, and it's a good fit for this kind of system. Can run in a cheap enbedded system, and it's easy to deploy and scale. 
The API should be documented with openAPI or similar, but due to time constraints this was not implemented.
The endpoints are:
-  /start a simple GET request to start the control loop, returns a simple status
-  /stop a simple GET request to start the control loop, returns a simple status
-  /temperature a get request to set the current temperature, the setpoint value shuld be passed as a query param, returns a simple status with the current value
-  /pressure a get request to set the current pressure, the setpoint value shuld be passed as a query param, returns a simple status with the current value
- /ws is the live websocket connection to the data stream, the data is sent as a json object with the current values of the sensors


### Error Handling

Errors are simply caught and presented to the user. The system is designed to be fault tolerant, and the user should be able to recover from any error.

### Testing Strategy

The system is designed to be tested at every level. The embedded system is tested using unit tests.  The logging system is tested using integration tests.
The UI and the API shpuld tested using integration tests, but due to the time constraints, this was not implemented.
