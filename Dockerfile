FROM golang

WORKDIR /app/testbench

ADD . /app

EXPOSE 7999

CMD ["go", "run", "interactivenode.go", "--port", "7999", "--interactive", "true"]
