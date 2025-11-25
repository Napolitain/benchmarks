use std::time::Instant;
use std::hint::black_box;

// Native Rust implementations
#[inline(never)]
fn fast_sum8(a: i64, b: i64, c: i64, d: i64, e: i64, f: i64, g: i64, h: i64) -> i64 {
    a + b + c + d + e + f + g + h
}

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
    println!("Rust Native Benchmark");
    println!("=====================\n");

    const FAST_ITERS: i64 = 1_000_000;
    const SLOW_ITERS: i32 = 100;
    const COMPUTE_ITERS: i32 = 1_000_000;

    // === FAST function ===
    println!("fast_sum8 (1M calls):");

    let start = Instant::now();
    let mut result: i64 = 0;
    for _ in 0..FAST_ITERS {
        result = black_box(fast_sum8(
            black_box(1), black_box(2), black_box(3), black_box(4),
            black_box(5), black_box(6), black_box(7), black_box(8)
        ));
    }
    let _ = black_box(result);
    let elapsed = start.elapsed();
    let total_ms = elapsed.as_secs_f64() * 1000.0;
    let per_call_ns = elapsed.as_nanos() as f64 / FAST_ITERS as f64;
    println!("  Total time:  {:8.2} ms", total_ms);
    println!("  Per call:    {:8.2} ns", per_call_ns);

    // === SLOW function ===
    println!("\nslow_compute (100 calls, 1M iters each):");

    let start = Instant::now();
    for i in 0..SLOW_ITERS {
        result = black_box(slow_compute(black_box(i as i64), black_box(COMPUTE_ITERS)));
    }
    let _ = black_box(result);
    let elapsed = start.elapsed();
    let total_ms = elapsed.as_secs_f64() * 1000.0;
    let per_call_ms = total_ms / SLOW_ITERS as f64;
    println!("  Total time:  {:8.2} ms", total_ms);
    println!("  Per call:    {:8.2} ms", per_call_ms);
}
