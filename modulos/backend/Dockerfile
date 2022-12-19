FROM golang:1.17

WORKDIR /Backend
COPY . .

RUN go get -u github.com/gorilla/mux
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get github.com/rs/cors
RUN go get github.com/go-playground/validator/v10
RUN go get github.com/joho/godotenv
RUN go get -d

EXPOSE 8001

CMD ["go", "run", "main.go"]
