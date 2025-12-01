use std::time::Instant;
use std::hint::black_box;

#[inline(never)]
fn slow_compute(seed: i64, iterations: i32) -> i64 {
    let mut h = seed as u64;
    for _ in 0..iterations {
        h ^= h >> 33;
        h = h.wrapping_mul(0xff51afd7ed558ccd);
        h ^= h >> 33;
        h = h.wrapping_mul(0xc4ceb9fe1a85ec53);
        h ^= h >> 33;
    }
    h as i64
}

fn main() {
    const SLOW_ITERS: i32 = 100;
    const COMPUTE_ITERS: i32 = 1_000_000;

    let start = Instant::now();
    let mut result: i64 = 0;
    for i in 0..SLOW_ITERS {
        result = black_box(slow_compute(black_box(i as i64), black_box(COMPUTE_ITERS)));
    }
    let _ = black_box(result);
    let elapsed = start.elapsed();
    
    let total_ms = elapsed.as_secs_f64() * 1000.0;
    let per_call_ms = total_ms / SLOW_ITERS as f64;
    println!("slow_compute: {:.2} ms total, {:.2} ms/call", total_ms, per_call_ms);
}
