IMAGE_NAME=forum

build:
	docker image build -t $(IMAGE_NAME) .

run:
	docker run -d -p 8081:8081 --name $(IMAGE_NAME)-container $(IMAGE_NAME)

stop:
	docker stop $(IMAGE_NAME)-container

remove:
	docker rm $(IMAGE_NAME)-container

clean: stop remove
	docker rmi $(IMAGE_NAME)

all: build run