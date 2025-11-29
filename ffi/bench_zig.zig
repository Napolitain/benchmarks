const std = @import("std");

// Native Zig implementations (no FFI, for comparison)
noinline fn fastSum8(a: i64, b: i64, c: i64, d: i64, e: i64, f: i64, g: i64, h: i64) i64 {
    return a + b + c + d + e + f + g + h;
}

noinline fn slowCompute(seed: i64, iterations: i32) i64 {
    var h: u64 = @bitCast(seed);
    for (0..@intCast(iterations)) |_| {
        h ^= h >> 33;
        h *%= 0xff51afd7ed558ccd;
        h ^= h >> 33;
        h *%= 0xc4ceb9fe1a85ec53;
        h ^= h >> 33;
    }
    return @bitCast(h);
}

fn doNotOptimize(val: anytype) @TypeOf(val) {
    asm volatile ("" : : [val] "r" (val) : "memory");
    return val;
}

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();

    try stdout.print("Zig Native Benchmark\n", .{});
    try stdout.print("====================\n\n", .{});

    const FAST_ITERS: usize = 1_000_000;
    const SLOW_ITERS: usize = 100;
    const COMPUTE_ITERS: i32 = 1_000_000;

    // === FAST function ===
    try stdout.print("fast_sum8 (1M calls):\n", .{});

    var timer = try std.time.Timer.start();
    var result: i64 = 0;
    for (0..FAST_ITERS) |_| {
        result = doNotOptimize(fastSum8(
            doNotOptimize(@as(i64, 1)),
            doNotOptimize(@as(i64, 2)),
            doNotOptimize(@as(i64, 3)),
            doNotOptimize(@as(i64, 4)),
            doNotOptimize(@as(i64, 5)),
            doNotOptimize(@as(i64, 6)),
            doNotOptimize(@as(i64, 7)),
            doNotOptimize(@as(i64, 8)),
        ));
    }
    _ = doNotOptimize(result);
    var elapsed = timer.read();
    var total_ms = @as(f64, @floatFromInt(elapsed)) / 1_000_000.0;
    const per_call_ns = @as(f64, @floatFromInt(elapsed)) / @as(f64, @floatFromInt(FAST_ITERS));

    try stdout.print("  Total time:  {d:8.2} ms\n", .{total_ms});
    try stdout.print("  Per call:    {d:8.2} ns\n", .{per_call_ns});

    // === SLOW function ===
    try stdout.print("\nslow_compute (100 calls, 1M iters each):\n", .{});

    timer = try std.time.Timer.start();
    for (0..SLOW_ITERS) |i| {
        result = doNotOptimize(slowCompute(doNotOptimize(@as(i64, @intCast(i))), doNotOptimize(COMPUTE_ITERS)));
    }
    _ = doNotOptimize(result);
    elapsed = timer.read();
    total_ms = @as(f64, @floatFromInt(elapsed)) / 1_000_000.0;
    const per_call_ms = total_ms / @as(f64, @floatFromInt(SLOW_ITERS));

    try stdout.print("  Total time:  {d:8.2} ms\n", .{total_ms});
    try stdout.print("  Per call:    {d:8.2} ms\n", .{per_call_ms});
}
