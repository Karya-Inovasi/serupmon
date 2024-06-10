FROM debian:buster-slim

RUN apt-get update && \
    apt-get install -y \
        build-essential \
        dh-make \
        wget \
        curl \
        git \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz \
    && rm go1.22.3.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /build

COPY . .

CMD ["make", "clean", "all"]
