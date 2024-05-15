FROM golang:1.21

WORKDIR /app 

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY .env .env


RUN CGO_ENABLED=0 go build -o /docker-go-crud 

EXPOSE 8080

CMD [ "/docker-go-crud" ]



# Reference : https://guptaaashutosh.hashnode.dev/how-to-build-docker-image-of-golang-project, https://medium.com/@guptaaashu354/how-to-create-docker-image-of-golang-project-7c01912afadb
# docker run -d -it –-rm -p [host_port]:[container_port] --name [container_name] [image_id/image_tag]

