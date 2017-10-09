BINARY := bin/jenkins-build-metrics
SOURCE := main.go

.PHONY: all clean

all: $(BINARY)

$(BINARY): $(SOURCE)
	CGO_ENABLED=0 go build -o $@ $<

clean:
	rm $(BINARY)
