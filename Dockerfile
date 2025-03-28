# syntax=docker/dockerfile:1

FROM golang:1.22.10

RUN apt-get update && apt-get install -y python3 python3-venv tree

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

#Copy Everything
COPY . .

#RUN ls
#RUN tree /

# Tax token detector
WORKDIR /app/apps/python/parsebytecode
RUN chmod +x setup.sh
RUN /bin/bash -c "./setup.sh"


WORKDIR /app/cmd

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o start /app/cmd


# Run
CMD ["/app/cmd/start"]