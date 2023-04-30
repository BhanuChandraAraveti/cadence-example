FROM golang:alpine

WORKDIR /app
COPY . .

RUN apk add --update make
RUN apk add supervisor



# RUN go mod download

# RUN ls

RUN go build -o bins/httpserver app/httpserver/main.go
RUN go build -o bins/worker app/worker/main.go

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Start Supervisor
CMD ["/usr/bin/supervisord"]

## root
#app
#
## current
#- app
#- bins