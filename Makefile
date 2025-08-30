.PHONY: build go run

USER_ID := 1
GROUP_ID := 1

ifneq ($(OS),Windows_NT)
  USER_ID := $(shell id -u)
  GROUP_ID := $(shell id -g)
endif

build:
	docker compose build

go:
	docker compose run --rm go bash -c "go $(filter-out $@,$(MAKECMDGOALS)) && chown -f -R $(USER_ID):$(GROUP_ID) . | true"

run:
	docker compose run --rm go bash -c "./$(filter-out $@,$(MAKECMDGOALS)) && chown -f -R $(USER_ID):$(GROUP_ID) . | true"
