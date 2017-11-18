docker:
	docker-compose build

up: docker
	docker-compose up

# silent
dockers:
	docker-compose build > build.out &

ups: dockers
	nohup docker-compose up > log.out &
stop:
	docker stop m0ney-server
	docker stop m0ney-db

go:
	go build .
	go build -o ./daemon/daemon ./daemon

gor: go
	./daemon/daemon &
	./m0ney