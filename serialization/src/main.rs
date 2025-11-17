use arrow::array::{ArrayRef, Int64Array, Float64Array, RecordBatch};
use arrow::datatypes::{DataType, Field, Schema};
use arrow::ipc::writer::StreamWriter;
use std::sync::Arc;
use std::time::Instant;

// Include generated protobuf code
pub mod benchmark {
    include!(concat!(env!("OUT_DIR"), "/benchmark.rs"));
}

use benchmark::BenchmarkData as ProtoBenchmarkData;

#[derive(Clone)]
struct TestRecord {
    id: i64,
    name: String,
    value: f64,
    numbers: Vec<i32>,
    data: Vec<u8>,
}

fn generate_test_data(size: usize) -> Vec<TestRecord> {
    (0..size)
        .map(|i| TestRecord {
            id: i as i64,
            name: format!("Record_{}", i),
            value: i as f64 * 1.5,
            numbers: vec![i as i32, (i + 1) as i32, (i + 2) as i32],
            data: vec![0u8; 64],
        })
        .collect()
}

fn benchmark_arrow(records: &[TestRecord]) -> Vec<u8> {
    let ids: Int64Array = records.iter().map(|r| Some(r.id)).collect();
    let values: Float64Array = records.iter().map(|r| Some(r.value)).collect();
    
    let schema = Schema::new(vec![
        Field::new("id", DataType::Int64, false),
        Field::new("value", DataType::Float64, false),
    ]);
    
    let batch = RecordBatch::try_new(
        Arc::new(schema.clone()),
        vec![
            Arc::new(ids) as ArrayRef,
            Arc::new(values) as ArrayRef,
        ],
    )
    .unwrap();
    
    let mut buffer = Vec::new();
    {
        let mut writer = StreamWriter::try_new(&mut buffer, &schema).unwrap();
        writer.write(&batch).unwrap();
        writer.finish().unwrap();
    }
    
    buffer
}

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

fn benchmark_capnproto(records: &[TestRecord]) -> Vec<u8> {
    let mut buffer = Vec::new();
    
    for record in records {
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

fn main() {
    println!("=== Serialization Benchmark Tool ===\n");
    
    let sizes = vec![100, 1000, 10000];
    
    for size in sizes {
        println!("Testing with {} records:", size);
        let records = generate_test_data(size);
        
        // Calculate approximate data size
        let data_size_mb = (size * (8 + 20 + 8 + 12 + 64)) as f64 / (1024.0 * 1024.0);
        
        // Benchmark Arrow
        let start = Instant::now();
        let arrow_data = benchmark_arrow(&records);
        let arrow_time = start.elapsed().as_secs_f64();
        let arrow_rate = data_size_mb / arrow_time;
        
        // Benchmark Protobuf
        let start = Instant::now();
        let protobuf_data = benchmark_protobuf(&records);
        let protobuf_time = start.elapsed().as_secs_f64();
        let protobuf_rate = data_size_mb / protobuf_time;
        
        // Benchmark Cap'n Proto
        let start = Instant::now();
        let capnproto_data = benchmark_capnproto(&records);
        let capnproto_time = start.elapsed().as_secs_f64();
        let capnproto_rate = data_size_mb / capnproto_time;
        
        println!("  Arrow:     {:.2} MB/s, {} bytes", arrow_rate, arrow_data.len());
        println!("  Protobuf:  {:.2} MB/s, {} bytes", protobuf_rate, protobuf_data.len());
        println!("  CapnProto: {:.2} MB/s, {} bytes", capnproto_rate, capnproto_data.len());
        println!();
    }
    
    println!("Run 'cargo bench' for detailed benchmarks using Criterion.");
}
