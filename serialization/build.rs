fn main() {
    // Compile protobuf schema
    prost_build::compile_protos(&["proto/benchmark.proto"], &["proto/"])
        .expect("Failed to compile protobuf schema");
    
    // Compile Cap'n Proto schema
    capnpc::CompilerCommand::new()
        .file("proto/benchmark.capnp")
        .run()
        .expect("Failed to compile Cap'n Proto schema");
}
