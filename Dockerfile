FROM golang:alpine AS build-stage

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN apk add --update make
RUN apk add supervisor


RUN GOOS=linux go build -o bins/httpserver app/httpserver/main.go
RUN GOOS=linux go build -o bins/worker app/worker/main.go

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY app/resources/application.yml /root/app/resources/application.yml

# RUN chmod +x bins/httpserver
# RUN chmod +x bins/worker

# Start a new stage from scratch
FROM alpine:latest AS build-release-stage

RUN apk --no-cache add ca-certificates supervisor
RUN apk add --no-cache bash
RUN apk --no-cache add tzdata

WORKDIR /app

COPY --from=build-stage /app/bins/httpserver bins/httpserver
COPY --from=build-stage /app/bins/worker bins/worker
COPY --from=build-stage /etc/supervisor/conf.d/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY --from=build-stage /root/app/resources/application.yml /root/app/resources/application.yml

# Expose port 3030 to the outside world
EXPOSE 3030

# Start Supervisor
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]