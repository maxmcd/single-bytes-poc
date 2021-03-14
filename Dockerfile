FROM golang

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH=$PATH:/root/.cargo/bin
RUN rustup target add wasm32-wasi

CMD cd /opt && cargo build --target wasm32-wasi && go run .
