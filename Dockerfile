FROM golang

WORKDIR /app

ADD . /app

EXPOSE 7999

CMD ["go", "run", "testbench/interactivenode.go", "--port", "7999", "--interactive", "true"]
