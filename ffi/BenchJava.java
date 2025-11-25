public class BenchJava {
    static final int FAST_ITERS = 1_000_000;
    static final int SLOW_ITERS = 100;
    static final int COMPUTE_ITERS = 1_000_000;

    // Native Java implementations
    static long fastSum8(long a, long b, long c, long d, long e, long f, long g, long h) {
        return a + b + c + d + e + f + g + h;
    }

    static long slowCompute(long seed, int iterations) {
        long h = seed;
        for (int i = 0; i < iterations; i++) {
            h ^= h >>> 33;
            h *= 0xff51afd7ed558ccdL;
            h ^= h >>> 33;
            h *= 0xc4ceb9fe1a85ec53L;
            h ^= h >>> 33;
        }
        return h;
    }

    // Volatile sink to prevent dead code elimination
    static volatile long sink;

    public static void main(String[] args) {
        System.out.println("Java Native Benchmark");
        System.out.println("=====================\n");

        // Warmup JIT
        for (int i = 0; i < 10000; i++) {
            sink = fastSum8(1, 2, 3, 4, 5, 6, 7, 8);
            sink = slowCompute(i, 1000);
        }

        // === FAST function ===
        System.out.println("fast_sum8 (1M calls):");

        long start = System.nanoTime();
        for (int i = 0; i < FAST_ITERS; i++) {
            sink = fastSum8(1, 2, 3, 4, 5, 6, 7, 8);
        }
        long end = System.nanoTime();

        double fastTotalMs = (end - start) / 1_000_000.0;
        double fastPerCallNs = (end - start) / (double) FAST_ITERS;

        System.out.printf("  Total time:  %8.2f ms%n", fastTotalMs);
        System.out.printf("  Per call:    %8.2f ns%n", fastPerCallNs);

        // === SLOW function ===
        System.out.println("\nslow_compute (100 calls, 1M iters each):");

        start = System.nanoTime();
        for (int i = 0; i < SLOW_ITERS; i++) {
            sink = slowCompute(i, COMPUTE_ITERS);
        }
        end = System.nanoTime();

        double slowTotalMs = (end - start) / 1_000_000.0;
        double slowPerCallMs = slowTotalMs / SLOW_ITERS;

        System.out.printf("  Total time:  %8.2f ms%n", slowTotalMs);
        System.out.printf("  Per call:    %8.2f ms%n", slowPerCallMs);
    }
}
