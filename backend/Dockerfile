# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /backend .

EXPOSE 8080

CMD [ "/backend" ]

# sudo docker build -t backend .
# sudo docker run -p 8080:8080 backend
# sudo docker run --network host --add-host=localhost:127.0.0.1 backend
# sudo docker run -d --network host --add-host=localhost:127.0.0.1 backend
# sudo docker ps
# sudo docker stop <container id>