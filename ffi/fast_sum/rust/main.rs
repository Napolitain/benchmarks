use std::time::Instant;
use std::hint::black_box;

#[inline(never)]
fn fast_sum8(a: i64, b: i64, c: i64, d: i64, e: i64, f: i64, g: i64, h: i64) -> i64 {
    a + b + c + d + e + f + g + h
}

fn main() {
    const FAST_ITERS: i64 = 1_000_000;

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
    println!("fast_sum8: {:.2} ms total, {:.2} ns/call", total_ms, per_call_ns);
}
