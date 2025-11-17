use criterion::{black_box, criterion_group, criterion_main, Criterion, BenchmarkId, Throughput};
use arrow::array::{ArrayRef, Int64Array, Float64Array};
use arrow::datatypes::{DataType, Field, Schema};
use arrow::record_batch::RecordBatch;
use arrow::ipc::writer::StreamWriter;
use std::sync::Arc;

// Include generated protobuf code
pub mod benchmark {
    include!(concat!(env!("OUT_DIR"), "/benchmark.rs"));
}

use benchmark::BenchmarkData as ProtoBenchmarkData;

// Generate test data
fn generate_test_data(size: usize) -> Vec<TestRecord> {
    (0..size)
        .map(|i| TestRecord {
            id: i as i64,
            name: format!("Record_{}", i),
            value: i as f64 * 1.5,
            numbers: vec![i as i32, (i + 1) as i32, (i + 2) as i32],
            data: vec![0u8; 64], // 64 bytes of data
        })
        .collect()
}

#[derive(Clone)]
struct TestRecord {
    id: i64,
    name: String,
    value: f64,
    numbers: Vec<i32>,
    data: Vec<u8>,
}

// Arrow serialization
fn benchmark_arrow(records: &[TestRecord]) -> Vec<u8> {
    // Create arrays
    let ids: Int64Array = records.iter().map(|r| Some(r.id)).collect();
    let values: Float64Array = records.iter().map(|r| Some(r.value)).collect();
    
    // Create schema
    let schema = Schema::new(vec![
        Field::new("id", DataType::Int64, false),
        Field::new("value", DataType::Float64, false),
    ]);
    
    // Create record batch
    let batch = RecordBatch::try_new(
        Arc::new(schema.clone()),
        vec![
            Arc::new(ids) as ArrayRef,
            Arc::new(values) as ArrayRef,
        ],
    )
    .unwrap();
    
    // Serialize
    let mut buffer = Vec::new();
    {
        let mut writer = StreamWriter::try_new(&mut buffer, &schema).unwrap();
        writer.write(&batch).unwrap();
        writer.finish().unwrap();
    }
    
    buffer
}

// Protobuf serialization
fn benchmark_protobuf(records: &[TestRecord]) -> Vec<u8> {
    use prost::Message;
    
    let mut buffer = Vec::new();
    for record in records {
        let proto_data = ProtoBenchmarkData {
            id: record.id,
            name: record.name.clone(),
            value: record.value,
            numbers: record.numbers.clone(),
            data: record.data.clone(),
        };
        
        proto_data.encode(&mut buffer).unwrap();
    }
    
    buffer
}

// Cap'n Proto serialization (simplified)
fn benchmark_capnproto(records: &[TestRecord]) -> Vec<u8> {
    // For simplicity, we'll use a basic binary format
    // A full Cap'n Proto implementation would require schema compilation
    let mut buffer = Vec::new();
    
    for record in records {
        // Simple binary encoding for demonstration
        buffer.extend_from_slice(&record.id.to_le_bytes());
        buffer.extend_from_slice(&(record.name.len() as u32).to_le_bytes());
        buffer.extend_from_slice(record.name.as_bytes());
        buffer.extend_from_slice(&record.value.to_le_bytes());
        buffer.extend_from_slice(&(record.numbers.len() as u32).to_le_bytes());
        for num in &record.numbers {
            buffer.extend_from_slice(&num.to_le_bytes());
        }
        buffer.extend_from_slice(&(record.data.len() as u32).to_le_bytes());
        buffer.extend_from_slice(&record.data);
    }
    
    buffer
}

fn encoding_benchmarks(c: &mut Criterion) {
    let sizes = vec![100, 1000, 10000];
    
    let mut group = c.benchmark_group("serialization_encoding");
    
    for size in sizes {
        let records = generate_test_data(size);
        let data_size = size * (8 + 20 + 8 + 12 + 64); // Approximate size per record
        
        group.throughput(Throughput::Bytes(data_size as u64));
        
        group.bench_with_input(BenchmarkId::new("Arrow", size), &records, |b, records| {
            b.iter(|| {
                let encoded = benchmark_arrow(black_box(records));
                black_box(encoded)
            });
        });
        
        group.bench_with_input(BenchmarkId::new("Protobuf", size), &records, |b, records| {
            b.iter(|| {
                let encoded = benchmark_protobuf(black_box(records));
                black_box(encoded)
            });
        });
        
        group.bench_with_input(BenchmarkId::new("CapnProto", size), &records, |b, records| {
            b.iter(|| {
                let encoded = benchmark_capnproto(black_box(records));
                black_box(encoded)
            });
        });
    }
    
    group.finish();
}

fn encoding_size_benchmark(_c: &mut Criterion) {
    let size = 1000;
    let records = generate_test_data(size);
    
    println!("\n=== Encoding Size Comparison (1000 records) ===");
    
    let arrow_size = benchmark_arrow(&records).len();
    println!("Arrow:     {} bytes", arrow_size);
    
    let protobuf_size = benchmark_protobuf(&records).len();
    println!("Protobuf:  {} bytes", protobuf_size);
    
    let capnproto_size = benchmark_capnproto(&records).len();
    println!("CapnProto: {} bytes", capnproto_size);
    
    println!("============================================\n");
}

criterion_group!(benches, encoding_benchmarks, encoding_size_benchmark);
criterion_main!(benches);
