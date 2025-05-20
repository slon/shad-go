FROM golang:1.24

RUN ["mkdir", "-p", "/build"]

WORKDIR /build

COPY . .

RUN ["go", "build", "-tags", "solution", "-o", "/firegod", "./firegod"]

ENTRYPOINT [ "/firegod", "-http=0.0.0.0:8080" ]
