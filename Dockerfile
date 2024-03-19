
FROM golang:1.21-alpine

ARG _SERVICE_NAME=similar_words_service
ARG _PROJECT_PATH="/${_SERVICE_NAME}"
ARG _EXECUTABLE_PATH="/server_executable/equivalent_words_server"

WORKDIR "${_PROJECT_PATH}"
COPY . .
RUN go get github.com/labstack/echo/v4


RUN echo go build -o "${_EXECUTABLE_PATH}"
RUN go build -o "${_EXECUTABLE_PATH}"
RUN chmod u+x "${_EXECUTABLE_PATH}"


# must hardcode the entrypoint, using the variable form will prevent outside signals from passing through
ENTRYPOINT [ "/server_executable/equivalent_words_server" ]