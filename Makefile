BINARY := bin/gomota

EXTERNAL_ASSETS := 
.PRECIOUS: ${EXTERNAL_ASSETS}

all: ${BINARY}

clean:
	rm -rf assets/external ${BINARY}

${BINARY}: ${EXTERNAL_ASSETS}
	go mod tidy
	go test
	go build -o $@

