const std = @import("std");

noinline fn fastSum8(a: i64, b: i64, c: i64, d: i64, e: i64, f: i64, g: i64, h: i64) i64 {
    return a + b + c + d + e + f + g + h;
}

fn doNotOptimize(val: anytype) @TypeOf(val) {
    asm volatile ("" : : [val] "r" (val) : "memory");
    return val;
}

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const FAST_ITERS: usize = 1_000_000;

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
    const elapsed = timer.read();
    
    const total_ms = @as(f64, @floatFromInt(elapsed)) / 1_000_000.0;
    const per_call_ns = @as(f64, @floatFromInt(elapsed)) / @as(f64, @floatFromInt(FAST_ITERS));
    try stdout.print("fast_sum8: {d:.2} ms total, {d:.2} ns/call\n", .{total_ms, per_call_ns});
}
