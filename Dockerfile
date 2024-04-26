# USAGE:
# docker build -t alum-challenges .
# docker run --publish 8000:8000 alum-challenges 

FROM golang:1.22 AS build-stage

WORKDIR /app

# copy to project directory and download deps
COPY . .
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o /alum-challenges

# Run tests in container
FROM build-stage AS run-test-stage
RUN go test -v ./..

# Deploy into a lean image
FROM gcr.io/distroless/base-debian11 AS run-stage
WORKDIR /app
COPY --from=build-stage /alum-challenges /alum-challenges
COPY . .
CMD ["tailwindcss-linux-x64" "-i" "static/input.css" "static/output.css" "--minify"]
EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT ["/alum-challenges"]
