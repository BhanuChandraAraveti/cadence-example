default: bins

httpserver:
	go build -o bins/httpserver app/httpserver/main.go

worker:
	go build -o bins/worker app/worker/main.go

bins: httpserver worker

clean:
	rm -rf bins