use criterion::{black_box, criterion_group, criterion_main, Criterion, BenchmarkId, Throughput};
use apache_avro::{Writer, Schema as AvroSchema};

// Include generated protobuf code
pub mod benchmark {
    include!(concat!(env!("OUT_DIR"), "/benchmark.rs"));
}

// Include generated Cap'n Proto code
pub mod benchmark_capnp {
    include!(concat!(env!("OUT_DIR"), "/proto/benchmark_capnp.rs"));
}

use benchmark::BenchmarkData as ProtoBenchmarkData;
use fory_derive::Fory;

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

#[derive(Clone, Fory)]
struct TestRecord {
    id: i64,
    name: String,
    value: f64,
    numbers: Vec<i32>,
    data: Vec<u8>,
}

// Avro serialization
fn benchmark_avro(records: &[TestRecord]) -> Vec<u8> {
    use apache_avro::types::{Record, Value};
    
    let schema_str = r#"
    {
      "type": "record",
      "name": "BenchmarkData",
      "namespace": "benchmark",
      "fields": [
        {"name": "id", "type": "long"},
        {"name": "name", "type": "string"},
        {"name": "value", "type": "double"},
        {"name": "numbers", "type": {"type": "array", "items": "int"}},
        {"name": "data", "type": "bytes"}
      ]
    }
    "#;
    
    let schema = AvroSchema::parse_str(schema_str).unwrap();
    let mut writer = Writer::new(&schema, Vec::new());
    
    for record in records {
        let mut avro_record = Record::new(&schema).unwrap();
        avro_record.put("id", Value::Long(record.id));
        avro_record.put("name", Value::String(record.name.clone()));
        avro_record.put("value", Value::Double(record.value));
        avro_record.put("numbers", Value::Array(
            record.numbers.iter().map(|&n| Value::Int(n)).collect()
        ));
        avro_record.put("data", Value::Bytes(record.data.clone()));
        
        writer.append(avro_record).unwrap();
    }
    
    writer.into_inner().unwrap()
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
    buffer
}

// Fory serialization
fn benchmark_fory(records: &[TestRecord]) -> Vec<u8> {
    let mut buffer = Vec::new();
    for record in records {
        let encoded = fory::serialize(record).unwrap();
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
        
        group.bench_with_input(BenchmarkId::new("Avro", size), &records, |b, records| {
            b.iter(|| {
                let encoded = benchmark_avro(black_box(records));
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

        group.bench_with_input(BenchmarkId::new("Fory", size), &records, |b, records| {
            b.iter(|| {
                let encoded = benchmark_fory(black_box(records));
                black_box(encoded)
            });
        });
    }
    
    group.finish();
}

criterion_group!(benches, encoding_benchmarks);
criterion_main!(benches);
