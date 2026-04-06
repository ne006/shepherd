BUILD_DIR = .build
APPS = $(notdir $(wildcard ./cmd/*))

build-%:
	go build -o $(BUILD_DIR)/$* ./cmd/$*

all: $(addprefix build-, $(APPS))

.PHONY: all build-%