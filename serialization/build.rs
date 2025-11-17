fn main() {
    // Compile protobuf schema
    prost_build::compile_protos(&["proto/benchmark.proto"], &["proto/"])
        .expect("Failed to compile protobuf schema");
    
    // Note: Cap'n Proto schemas are typically compiled separately
    // For this benchmark, we'll use runtime compilation
}
