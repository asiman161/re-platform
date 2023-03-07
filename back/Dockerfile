FROM rust:1.67.0 as build

RUN USER=root cargo new --bin back
WORKDIR /back

COPY ./Cargo.lock ./Cargo.lock
COPY ./Cargo.toml ./Cargo.toml

# 3. Build only the dependencies to cache them
RUN cargo build --release
RUN rm src/*.rs

COPY ./src ./src
RUN rm ./target/release/deps/back*
RUN cargo build --release

COPY ./src ./src

# 5. Build for release.
RUN rm ./target/release/deps/back*
#RUN cargo install --path .
RUN cargo build --release

CMD ["./target/release/back"]
