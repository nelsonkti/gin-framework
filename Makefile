.PHONY: start stop restart debug
SHELL := /bin/bash
#打包
build:
	cd main && go mod tidy && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(service) $(service).go
	

# Variables
NS=ms
POD_NAME_SRC=$(filter-out $@,$(MAKECMDGOALS))

# Find pod name based on input
find-pod-name:
	@echo "After if statement"
	$(eval POD_NAME_SRC=$(shell echo $(POD_NAME_SRC) | awk '{print $$2}' ))
	@echo "POD_NAME_SRC: [$(POD_NAME_SRC)]"
	$(eval POD2_NAME=$(shell kubectl -n $(NS) get pod | grep "$(POD_NAME_SRC)" | grep "Running" | awk '{print $$1}'  || echo "NO_POD_FOUND"))
	@echo "POD2_NAME after kubectl command: [$(POD2_NAME)]"
	$(eval POD2_NAME=$(shell kubectl -n $(NS) get pod | awk '{print $$1}' | grep "$(POD_NAME_SRC)-"))
#重启
restart: find-pod-name
	#-kubectl -n $(NS) exec -it $(POD2_NAME) -- pkill -9 /app/main/$(service)
 	#kubectl -n $(NS) exec -it $(POD2_NAME) -- nohup /app/main/$(service) &
	kubectl rollout restart deployment/$(POD_NAME_SRC) -n $(NS)
#启动
start: find-pod-name
	kubectl -n $(NS) exec -it $(POD2_NAME) -- nohup /app/main/$(service) &
#停止
stop: find-pod-name
	kubectl -n $(NS) exec -it $(POD2_NAME) -- pkill -9 /app/main/$(service)
#状态
status: find-pod-name
	kubectl -n $(NS) exec -it $(POD2_NAME) -- ps -ef
#打印运行日志
debug: find-pod-name
	@if [ "$(POD2_NAME)" = "NO_POD_FOUND" ]; then \
		echo "No pod found, cannot start service"; \
	else \
		kubectl -n $(NS) exec -it $(POD2_NAME) -- /app/main/$(service); \
	fi
