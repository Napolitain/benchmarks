use criterion::{black_box, criterion_group, criterion_main, Criterion, BenchmarkId, Throughput};
use arrow::array::{ArrayRef, Int64Array, Float64Array, StringArray, BinaryArray, Int32Array};
use arrow::datatypes::{DataType, Field, Schema};
use arrow::record_batch::RecordBatch;
use arrow::ipc::writer::StreamWriter;
use std::sync::Arc;

// Include generated protobuf code
pub mod benchmark {
    include!(concat!(env!("OUT_DIR"), "/benchmark.rs"));
}

// Include generated Cap'n Proto code
pub mod benchmark_capnp {
    include!(concat!(env!("OUT_DIR"), "/proto/benchmark_capnp.rs"));
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
    // Create arrays for all fields
    let ids: Int64Array = records.iter().map(|r| Some(r.id)).collect();
    let names: StringArray = records.iter().map(|r| Some(r.name.as_str())).collect();
    let values: Float64Array = records.iter().map(|r| Some(r.value)).collect();
    
    // Create ListArray for numbers field
    let numbers_field = Arc::new(Field::new("item", DataType::Int32, false));
    let mut numbers_builder = arrow::array::ListBuilder::new(Int32Array::builder(0)).with_field(numbers_field);
    for record in records {
        let int_builder = numbers_builder.values();
        for &num in &record.numbers {
            int_builder.append_value(num);
        }
        numbers_builder.append(true);
    }
    let numbers_array = numbers_builder.finish();
    
    // Create BinaryArray for data field
    let data_array: BinaryArray = records.iter().map(|r| Some(r.data.as_slice())).collect();
    
    // Create schema
    let schema = Schema::new(vec![
        Field::new("id", DataType::Int64, false),
        Field::new("name", DataType::Utf8, false),
        Field::new("value", DataType::Float64, false),
        Field::new("numbers", DataType::List(Arc::new(Field::new("item", DataType::Int32, false))), false),
        Field::new("data", DataType::Binary, false),
    ]);
    
    // Create record batch
    let batch = RecordBatch::try_new(
        Arc::new(schema.clone()),
        vec![
            Arc::new(ids) as ArrayRef,
            Arc::new(names) as ArrayRef,
            Arc::new(values) as ArrayRef,
            Arc::new(numbers_array) as ArrayRef,
            Arc::new(data_array) as ArrayRef,
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

// Cap'n Proto serialization
fn benchmark_capnproto(records: &[TestRecord]) -> Vec<u8> {
    use capnp::message::Builder;
    use capnp::serialize;
    
    let mut buffer = Vec::new();
    
    for record in records {
        let mut message = Builder::new_default();
        {
            let mut benchmark_data = message.init_root::<benchmark_capnp::benchmark_data::Builder>();
            benchmark_data.set_id(record.id);
            benchmark_data.set_name(&record.name);
            benchmark_data.set_value(record.value);
            
            let mut numbers_builder = benchmark_data.reborrow().init_numbers(record.numbers.len() as u32);
            for (i, &num) in record.numbers.iter().enumerate() {
                numbers_builder.set(i as u32, num);
            }
            
            benchmark_data.set_data(&record.data);
        }
        
        let encoded = serialize::write_message_to_words(&message);
        buffer.extend_from_slice(&encoded);
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

criterion_group!(benches, encoding_benchmarks);
criterion_main!(benches);
