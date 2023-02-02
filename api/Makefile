VER ?= 4

run:
	go run main.go

docker-build:
	docker build -t fruits:$(VER) -f refs/Dockerfile.$(VER) .

docker-run:
	docker run -it -d --name fruits -p 9999:9999 fruits:$(VER)

docker-rm:
	docker rm -f fruits

zip:
	git archive --format zip --output fruits.zip master

devbox-setup:
	sudo bash setup.sh

demo:
	curl "localhost:9999/"
	curl "localhost:9999/buy?fruit=apple&count=10"
	curl "localhost:9999/buy?fruit=orange&count=20"
	curl "localhost:9999/sell?fruit=apple&count=5"
	curl "localhost:9999/sell?fruit=orange&count=10"
	curl "localhost:9999/"
