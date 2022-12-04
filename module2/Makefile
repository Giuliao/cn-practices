build:

.PHONY: build

run:
	VERSION=2 go run main.go -logtostderr=true -port=8080 -v=4
.PHONY: run

build-image:
	docker build -f Dockerfile . --target build --tag build-stage-image
	docker build -f Dockerfile .  -t giuliao/module3-httpserver:v3.0
	docker rmi build-stage-image # 删除中间产物
.PHONY: run