build:

.PHONY: build

run:
	VERSION=2 go run main.go -logtostderr=true -port=8080
.PHONY: run

build-image:
	docker build -f Dockerfile . --target build --tag build-stage-image
	docker build -f Dockerfile .  -t giuliao/module3-httpserver:v2.1
	docker rmi build-stage-image # 删除中间产物
.PHONY: run