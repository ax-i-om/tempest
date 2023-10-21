FROM golang:1.20-alpine AS BUILD

WORKDIR /app

# download the required Go dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download
#COPY *.go ./
COPY . ./

RUN go build -o /tempest

EXPOSE 8080

CMD [ "/tempest", "json", "testing"]