.PHONY: stop-all
stop-all:
	-if [ -f authsvc.pid ]; then kill `cat authsvc.pid` 2>/dev/null || true; rm -f authsvc.pid; fi
	-if [ -f statesvc.pid ]; then kill `cat statesvc.pid` 2>/dev/null || true; rm -f statesvc.pid; fi
	docker-compose down || true

.PHONY: build-start
build-start:
	docker-compose up -d
	sh build.sh
	@if [ ! -f authentication/authsvc ] || [ ! -f statesvc/statesvc ]; then echo "Build failed: binaries not created"; exit 1; fi
	cd authentication && ./authsvc > authsvc.log 2>&1 & echo $$! > authsvc.pid
	cd statesvc && ./statesvc > statesvc.log 2>&1 & echo $$! > statesvc.pid
	cd frontend && npm start 

.PHONY: start
start:
	@if [ ! -f authentication/authsvc ]; then echo "authsvc not found, run 'make build-start' first"; exit 1; fi
	@if [ ! -f statesvc/statesvc ]; then echo "statesvc not found, run 'make build-start' first"; exit 1; fi
	cd authentication && ./authsvc > authsvc.log 2>&1 & echo $$! > authsvc.pid
	cd statesvc && ./statesvc > statesvc.log 2>&1 & echo $$! > statesvc.pid

.PHONY: clean
clean:
	rm -f authentication/authsvc authentication/authsvc.log
	rm -f statesvc/statesvc statesvc/statesvc.log
	rm -f authsvc.pid statesvc.pid
	docker-compose down -v --remove-orphans || true

.PHONY: build-docker
build-docker:
	sh build.sh
	docker build -t authsvc -f authentication/dockerfile
	docker build -t statesvc -f statesvc/dockerfile
