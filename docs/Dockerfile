FROM alpine:latest
RUN apk add --no-cache wget
WORKDIR /docs
RUN wget https://github.com/gohugoio/hugo/releases/download/v0.53/hugo_0.53_Linux-64bit.tar.gz && \
  mkdir hugo && \
  tar -zxvf hugo_0.53_Linux-64bit.tar.gz -C hugo/ && \
  mv hugo/hugo /usr/local/bin/hugo && \
  rm -rf hugo*
COPY . .
RUN hugo
