const std = @import("std");

const RectangleData = struct {
    a: f64,
    b: f64,
    c: f64,
    d: f64,
};

fn parseSimpleYaml(contents: []const u8) !RectangleData {
    var data = RectangleData{ .a = 0, .b = 0, .c = 0, .d = 0 };
    var lines = std.mem.splitScalar(u8, contents, '\n');

    while (lines.next()) |line| {
        const trimmed = std.mem.trim(u8, line, " \t\r");
        if (trimmed.len == 0) continue;

        if (std.mem.indexOfScalar(u8, trimmed, ':')) |colon_idx| {
            const key = std.mem.trim(u8, trimmed[0..colon_idx], " \t");
            const value_str = std.mem.trim(u8, trimmed[colon_idx + 1 ..], " \t");
            const value = std.fmt.parseFloat(f64, value_str) catch continue;

            if (std.mem.eql(u8, key, "a")) {
                data.a = value;
            } else if (std.mem.eql(u8, key, "b")) {
                data.b = value;
            } else if (std.mem.eql(u8, key, "c")) {
                data.c = value;
            } else if (std.mem.eql(u8, key, "d")) {
                data.d = value;
            }
        }
    }
    return data;
}

fn computeRectangleArea(data: RectangleData) f64 {
    const width = @abs(data.c - data.a);
    const height = @abs(data.d - data.b);
    return width * height;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const args = try std.process.argsAlloc(allocator);
    defer std.process.argsFree(allocator, args);

    if (args.len < 2) {
        std.debug.print("Usage: rectangle <yaml-file>\n", .{});
        return;
    }

    const yaml_file = args[1];

    var timer = try std.time.Timer.start();

    const file_contents = try std.fs.cwd().readFileAlloc(allocator, yaml_file, 1024 * 1024);
    defer allocator.free(file_contents);

    const data = try parseSimpleYaml(file_contents);

    const area = computeRectangleArea(data);

    const elapsed = timer.read();
    const elapsed_ms = @as(f64, @floatFromInt(elapsed)) / 1_000_000.0;

    const stdout = std.io.getStdOut().writer();
    try stdout.print("Rectangle area: {d:.2}\n", .{area});
    try stdout.print("Time: {d:.6} ms\n", .{elapsed_ms});
}
