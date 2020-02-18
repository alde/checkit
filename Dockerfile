FROM library/golang:1.13-alpine

COPY go.mod go.sum /app/
RUN cd /app && go mod download
COPY checkstyle/*.go /app/checkstyle/
COPY spotbugs/*.go /app/spotbugs/
COPY *.go /app/
RUN cd /app && go build .

ENTRYPOINT [ "/app/checkit" ]