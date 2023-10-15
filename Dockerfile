FROM golang:1.21.1
WORKDIR /app
COPY . .
ENV BASE_DIR="/data"
ENV RESULT_DIR="result"
ENV ICONS_DIR="icons"
RUN ["go", "build", "-o", "app"]
EXPOSE 8081
CMD ["./app"]
