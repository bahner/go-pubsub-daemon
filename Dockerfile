ARG BUILD_IMAGE=golang:1.19.11-alpine
#Dockerfile
#From which image we want to build. This is basically our environment.
FROM ${BUILD_IMAGE} as Build

#This will copy all the files in our repo to the inside the container at root location.
COPY . . 

#build our binary at root location. Binary name will be main. We are using go modules so gpath env variable should be empty.
RUN GOPATH= go build -o /main main.go 

FROM alpine:latest
COPY --from=Build /main /main

#we tell docker what to run when this image is run and run it as executable.
ENTRYPOINT [ "/main" ]