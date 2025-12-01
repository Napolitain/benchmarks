const std = @import("std");

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
    const SLOW_ITERS: usize = 100;
    const COMPUTE_ITERS: i32 = 1_000_000;

    var timer = try std.time.Timer.start();
    var result: i64 = 0;
    for (0..SLOW_ITERS) |i| {
        result = doNotOptimize(slowCompute(doNotOptimize(@as(i64, @intCast(i))), doNotOptimize(COMPUTE_ITERS)));
    }
    _ = doNotOptimize(result);
    const elapsed = timer.read();
    
    const total_ms = @as(f64, @floatFromInt(elapsed)) / 1_000_000.0;
    const per_call_ms = total_ms / @as(f64, @floatFromInt(SLOW_ITERS));
    try stdout.print("slow_compute: {d:.2} ms total, {d:.2} ms/call\n", .{total_ms, per_call_ms});
}
