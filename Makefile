EXECUTABLE=go-hellokube
DOCKERIMAGE=brooksy/go-hellokube
DOCKERTAG=v0.01

# Determine architecture
BUILDARCH 	:=
ifeq ($(OS),Windows_NT)
    FILEEXTENSION = .exe
    GOOS = windows
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
    	GOARCH = amd64
    else
    	GOARCH = 386
    endif

	BIN_DIR=.\bin
	CMDQUIET = >nul 2>nul & verify >nul

	CLEAN_ALL = rmdir $(BIN_DIR)\ /s /q && mkdir $(BIN_DIR)
	RM_DIR = rmdir /q /s $(BIN_DIR)\$(GOOS)$(GOARCH)
	RM_EXEC = del $(BIN_DIR)\$(EXECUTABLE)$(FILEEXTENSION)
	CP_EXEC = copy $(BIN_DIR)\$(GOOS)$(GOARCH)\$(EXECUTABLE)$(FILEEXTENSION) $(BIN_DIR)\$(EXECUTABLE)$(FILEEXTENSION)
	CP_LINUX_AMD64 = copy $(BIN_DIR)\linuxamd64\$(EXECUTABLE) $(BIN_DIR)\$(EXECUTABLE)

	# Create the compile statements for each OS	and architecture type
	FREEBSD_AMD64 = cmd /V /C "set GOOS=freebsd&&set GOARCH=amd64&& go build -o $(BIN_DIR)/freebsdamd64/$(EXECUTABLE)"
	LINUX_AMD64 = cmd /V /C "set GOOS=linux&&set GOARCH=amd64&& go build -o $(BIN_DIR)/linuxamd64/$(EXECUTABLE)"
	LINUX_ARM64 = cmd /V /C "set GOOS=linux&&set GOARCH=arm64&& go build -o $(BIN_DIR)/linuxarm64/$(EXECUTABLE)"
	OSX_AMD64 = cmd /V /C "set GOOS=darwin&&set GOARCH=amd64&& go build -o $(BIN_DIR)/darwinamd64/$(EXECUTABLE)"
	OSX_ARM64 = cmd /V /C "set GOOS=darwin&&set GOARCH=arm64&& go build -o $(BIN_DIR)/darwinarm64/$(EXECUTABLE)"
	WIN_ARM64 = cmd /V /C "set GOOS=windows&&set GOARCH=arm64&& go build -o $(BIN_DIR)/windowsarm64/$(EXECUTABLE).exe"
	WIN_AMD64 = cmd /V /C "set GOOS=windows&&set GOARCH=amd64&& go build -o $(BIN_DIR)/windowsamd64/$(EXECUTABLE).exe"
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS = linux
	endif
	ifeq ($(UNAME_S),Darwin)
		GOOS = darwin
	endif
	UNAME_P := $(shell uname -p)
	ifeq ($(UNAME_P),x86_64)
		GOARCH = amd64
	endif
	ifneq ($(filter %86,$(UNAME_P)),)
    	GOARCH = 386
	endif
	ifneq ($(filter arm%,$(UNAME_P)),)
    	GOARCH = arm64
	endif

	BIN_DIR=./bin
	CMDQUIET= >/dev/null 2>&1

	CLEAN_ALL = rm -rf $(BIN_DIR)/*
	RM_DIR = rm -rf $(BIN_DIR)/$(GOOS)$(GOARCH)
	RM_EXEC = rm -f $(BIN_DIR)/$(EXECUTABLE)$(FILEEXTENSION)
	CP_EXEC = cp $(BIN_DIR)/$(GOOS)$(GOARCH)/$(EXECUTABLE)$(FILEEXTENSION) $(BIN_DIR)/$(EXECUTABLE)$(FILEEXTENSION)
	CP_LINUX_AMD64 = cp $(BIN_DIR)/linuxamd64/$(EXECUTABLE) $(BIN_DIR)/$(EXECUTABLE)

	# Create the compile statements for each OS	and architecture type
	FREEBSD_AMD64 = GOOS=freebsd GOARCH=amd64 go build -o $(BIN_DIR)/freebsdamd64/$(EXECUTABLE)
	LINUX_AMD64 = GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/linuxamd64/$(EXECUTABLE)
	LINUX_ARM64 = GOOS=linux GOARCH=arm64 go build -o $(BIN_DIR)/linuxarm64/$(EXECUTABLE)
	OSX_AMD64 = GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/darwinamd64/$(EXECUTABLE)
	OSX_ARM64 = GOOS=darwin GOARCH=arm64 go build -o $(BIN_DIR)/darwinarm64/$(EXECUTABLE)
	WIN_AMD64 = GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/windowsamd64/$(EXECUTABLE).exe
	WIN_ARM64 = GOOS=windows GOARCH=arm64 go build -o $(BIN_DIR)/windowsarm64/$(EXECUTABLE).exe
endif

all: clean build run

run:
	$(BIN_DIR)/$(EXECUTABLE)$(FILEEXTENSION)

docker:
	$(LINUX_AMD64)
	$(CP_LINUX_AMD64)
	docker build -t $(DOCKERIMAGE):$(DOCKERTAG) -f Dockerfile .

k8sup:
	kubectl create -f ./k8s-deployment.yaml

k8sdown:
	kubectl delete deployment $(EXECUTABLE)
	kubectl delete service $(EXECUTABLE)
	kubectl delete ingress $(EXECUTABLE)

build:
	@echo Compile for GOOS: $(GOOS), GOARCH: $(GOARCH), Extension: $(FILEEXTENSION)
# build with debug off to get a smaller executable
# go build -ldflags '-w -s' -o $(BIN_DIR)/$(GOOS)$(GOARCH)/$(EXECUTABLE)$(FILEEXTENSION)
	go build  -o $(BIN_DIR)/$(GOOS)$(GOARCH)/$(EXECUTABLE)$(FILEEXTENSION)
	$(CP_EXEC)

clean:
	-$(RM_DIR)
	-$(RM_EXEC)

build_all:
	$(FREEBSD_AMD64)
	$(LINUX_AMD64)
	$(LINUX_ARM64)
	$(OSX_AMD64)
	$(OSX_ARM64)
	$(WIN_AMD64)
	$(WIN_ARM64)

clean_all:
	@$(CLEAN_ALL)



