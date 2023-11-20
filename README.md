# json-go-mimic

## Description
`json-go-mimic` is a simple HTTP server written in Go, designed to mimic JSON responses for testing and development purposes. It allows for easy configuration of endpoints to serve custom JSON responses and supports various methods of HTTP authentication.

## Features
- Easy configuration for multiple endpoints.
- Customizable JSON response for each endpoint.
- Support for Bearer Token and API Key based authentication.
- Docker support for simple deployment and scaling.

## Installation and Usage
### Running Locally
Make sure you have the latest version of Go installed on your system. You can install Go using your package manager or by downloading it from the [official Go website](https://golang.org/dl/).

To run `json-go-mimic` locally, clone the repository and start the server:
```bash
git clone https://github.com/erolatex/json-go-mimic.git
cd json-go-mimic
go run src/main.go
```

### Running with Docker
Build and run the `json-go-mimic` server using Docker with the following commands:
```bash
docker build -t json-go-mimic .
docker run -d -p 7732:7732 -v $(pwd)/configs:/app/configs -v $(pwd)/data:/app/data json-go-mimic
```

## Configuration
Edit the `configs/config.json` file to set up your endpoints. An example configuration file would look like this, with your specific endpoints and settings:
```json
{
  "port": 7732,
  "endpoints": [
    {
      "path": "/endpoint/bearer",
      "jsonFilePath": "data/first.json",
      "authType": "Bearer",
      "authKey": "your_bearer_token_here"
    },
    {
      "path": "/endpoint/apikey",
      "jsonFilePath": "data/second.json",
      "authType": "X-Api-Key",
      "authKey": "your_api_key_here"
    },
    {
      "path": "/public/endpoint",
      "jsonFilePath": "data/third.json",
      "authType": "None",
      "authKey": ""
    }
  ]
}
```

## License
This project is licensed under the BSD-3 License.

## Author
[erolatex](https://github.com/erolatex)
