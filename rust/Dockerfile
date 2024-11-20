FROM rust:latest as builder

WORKDIR /text-similarity-rust

COPY . .
RUN cargo build --release

FROM gcr.io/distroless/cc-debian12

WORKDIR /app
COPY --from=builder /text-similarity-rust/target/release/text-similarity-rust .

CMD ["./text-similarity-rust"]