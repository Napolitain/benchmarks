// Library for serialization benchmarks
// This crate only provides benchmarks via `cargo bench`

// Include generated protobuf code
pub mod benchmark {
    include!(concat!(env!("OUT_DIR"), "/benchmark.rs"));
}

// Include generated Cap'n Proto code
pub mod benchmark_capnp {
    include!(concat!(env!("OUT_DIR"), "/proto/benchmark_capnp.rs"));
}
