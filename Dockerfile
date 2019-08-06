FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Nick Hong <nick.hong@klook.com>"

RUN mkdir /app
# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN export GO15VENDOREXPERIMENT=1 && export GO111MODULE=on && export GOPROXY=https://goproxy.io && go build

# Expose port 8080 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./gwyApp"]