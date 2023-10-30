# Start from the latest golang base image
FROM golang:1.21.3

LABEL author="Mauricio Milano"
LABEL description="jobsity challenge"
LABEL version="1.0"


WORKDIR /worker/
# COPY go.mod ./go.mod
COPY go.sum go.sum ./ 

# Copy the source from the current directory to the Working Directory inside the container
COPY . .


ENV PUBLISHER_QUEUE='pqueue'
ENV RECEIVER_QUEUE='rqueue'
ENV RMQ_HOST='rabbitmq'
ENV RMQ_USERNAME='admin'
ENV RMQ_PASSWORD='admin'
ENV RMQ_PORT='5672'



RUN go mod download



CMD ["go","run","main.go"]