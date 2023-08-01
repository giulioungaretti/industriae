
# API
This is a Go project that uses a dev container for development.

### Prerequisites
Docker
Visual Studio Code
Remote - Containers extension

## Getting Started -- DEV 
Clone this repository
Open the repository in Visual Studio Code.
Install the recommended extensions.

###  DEV 
Reopen the repository in a dev container by clicking the "Reopen in Container" button in the bottom right corner of the window.
Run the project by opening a terminal and running go run main.go.
Test the project by opening a terminal and running go test ./...

### Build
Use the included docker file to build the API server.
The port is harcoded in this example, it should be an env variable in prod!
The server accepts CORS requests from localhost:3000, this should be changed in prod!