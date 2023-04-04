ifeq ($(OS),Windows_NT)
  EXEEXT = .exe
endif

build:
	go build -o ./replacer$(EXEEXT) main.go

build-static:
	go build --ldflags '-linkmode external -extldflags=-static' -o ./replacer$(EXEEXT) ./main.go
