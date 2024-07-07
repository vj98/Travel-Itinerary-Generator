## To run the app use
`go run main.go`

## To build the binary file
`GOOS=linux GOARCH=amd64 go build -o bootstrap main.go`

## zip file to upload on lambda
`zip bootstrap.zip bootstrap PROD.env`

## TO execute docker file
 `docker build -t my-go-app .`

## run application in docker
`docker run -p 8080:8080 my-go-app`