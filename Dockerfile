FROM golang:latest AS builder
WORKDIR /app
COPY ./ /app/


RUN go mod tidy 
RUN go mod download 
RUN make build


FROM gcr.io/distroless/cc-debian11
WORKDIR /
# run by default run sh and is not available on distroless Workaround
#COPY /db /db
COPY --from=builder /app/assets /assets
#COPY /posts /posts
COPY --from=builder /app/bin/blog-me /blog-me
CMD [ "/blog-me" ]