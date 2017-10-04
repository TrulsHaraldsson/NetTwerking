FROM golang

WORKDIR /app

ADD . /app

EXPOSE 7999

# boolean, true if it is an interactive node.
ARG interactive

RUN echo "interactive is $interactive"

CMD ["go", "run", "testbench/interactivenode.go", "--port", "7999", "--interactive", "$interactive"]
