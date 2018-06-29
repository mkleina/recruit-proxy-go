BUILD_DIR=$(shell pwd)

build:
	go build -o server github.com/go-service/proxy_server
	go build -o client github.com/go-service/proxy_client
swagger:
	swagger generate server -t proxy_server -f swagger.yml --exclude-main -A "Recruit proxy server"
	swagger generate client -t proxy_client -f swagger.yml -A "Recruit proxy client"
deps:
	curl https://glide.sh/get | sh
	go get github.com/go-swagger/go-swagger/cmd/swagger
	cd $(BUILD_DIR)/proxy_server && glide install
	cd $(BUILD_DIR)/proxy_client && glide install
clean:
	rm -f $(BUILD_DIR)/server
	rm -f $(BUILD_DIR)/client
	rm -rf $(BUILD_DIR)/proxy_server/vendor
	rm -rf $(BUILD_DIR)/proxy_client/vendor


