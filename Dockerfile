FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Nick Hong <nick.hong@klook.com>"

RUN mkdir /app
RUN export GO111MODULE=on
RUN export GOPROXY=https://goproxy.io
# We copy everything in the root directory
# into our /app directory
ADD . /app
# We specify that we now wish to execute
# any further commands inside our /app
# directory
WORKDIR /app
# we run go build to compile the binary
# executable of our Go program
RUN go build -o main .
# Expose port 8080 to the outside world
EXPOSE 8000
# Our start command which kicks off
# our newly created binary executable
CMD ["/app/main"]