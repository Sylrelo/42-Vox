GOSRC = $(shell ls -1 *.go)

install: 
	@printf "Installing dependencies...\n"
	@go get ./...

run:
	@go run -gcflags '-dwarf=0 -B -wb=0 -C=4' -ldflags '-s -w' -tags tuneparam $(GOSRC)

race:
	@go run -race $(GOSRC) -fs=false

build:
	@go build -gcflags '-dwarf=0 -B -wb=0 -C=4' -ldflags '-s -w' -tags tuneparam
	@./vox -help

clean:
	@rm ./vox