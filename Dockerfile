FROM golang:alpine as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git 

WORKDIR /app

COPY go.* ./ 
RUN go mod download

COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . 


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env . 

ENV GIN_MODE release
EXPOSE 8000

CMD ["./main"]



