# Stage Build
FROM golang:1.11 as build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
COPY *.go ./
RUN go mod download
RUN go build


# Stage Run
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/gtbchart .
RUN chmod +x gtbchart
CMD ["./gtbchart"]
EXPOSE 7000