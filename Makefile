tag = latest
env = staging
name = osc
port = 10001

.PHONY: build

deploy:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/counter ./cmd/counter 
	docker build -t $(name):$(tag) .
	uglifyjs --compress dead_code,evaluate,booleans,loops,unused,hoist_funs,hoist_vars,if_return,join_vars --mangle -o web/script.min.js web/script.js
	# cp web/index.html web/*.min.js /opt/$(name)/web
	docker stop $(name)_$(env); exit 0
	docker rm $(name)_$(env); exit 0
	docker run --name $(name)_$(env) -p$(port):$(port) $(name):$(tag)
