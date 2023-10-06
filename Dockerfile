FROM golang:1.20-alpine as builder

WORKDIR /backend/

COPY go.* ./ 
RUN go mod download

COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /backend/bin/app .


FROM alpine:latest

WORKDIR /backend/

COPY --from=builder /backend/bin/app /backend/bin/app
COPY .env /backend/.env

ENV GIN_MODE release
EXPOSE 8000

RUN chmod +x /backend/bin/app 
ENTRYPOINT [ "/backend/bin/app" ]



