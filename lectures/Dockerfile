FROM golang:1.13

RUN CGO_ENABLED=0 go get golang.org/x/tools/cmd/present && go install golang.org/x/tools/cmd/present

COPY . /lectures

WORKDIR /lectures

ENTRYPOINT /go/bin/present -http ':3999' -orighost p.go.manytask.org -use_playground -play=0
