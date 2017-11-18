docker:
	docker-compose build

up: docker
	docker-compose up

# silent
dockers:
	docker-compose build > build.out &

ups: dockers
	nohup docker-compose up > log.out &


go:
	go build .
	go build -o ./daemon/daemon ./daemon

gor: go
	./daemon/daemon &
	./m0ney