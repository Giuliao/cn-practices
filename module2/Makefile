build:

.PHONY: build

run:
	go run main.go -logtostderr true
.PHONY: run


build-image:
	docker build -f Dockerfile . --target build --tag module2-build-stage
	docker build -f Dockerfile .  -t giuliao/module2-demo:v1.0
	docker rmi module2-build-stage # 删除中间产物
.PHONY: run