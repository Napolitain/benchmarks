const FAST_ITERS = 1_000_000;
const SLOW_ITERS = 100;
const COMPUTE_ITERS = 1_000_000;

// Native TypeScript implementations
function fastSum8(a: bigint, b: bigint, c: bigint, d: bigint, e: bigint, f: bigint, g: bigint, h: bigint): bigint {
    return a + b + c + d + e + f + g + h;
}

function slowCompute(seed: bigint, iterations: number): bigint {
    let h = BigInt.asUintN(64, seed);
    for (let i = 0; i < iterations; i++) {
        h ^= h >> 33n;
        h = BigInt.asUintN(64, h * 0xff51afd7ed558ccdn);
        h ^= h >> 33n;
        h = BigInt.asUintN(64, h * 0xc4ceb9fe1a85ec53n);
        h ^= h >> 33n;
    }
    return h;
}

function main() {
    console.log("Node.js/TypeScript Native Benchmark");
    console.log("====================================\n");

    // === FAST function ===
    console.log("fast_sum8 (1M calls):");

    let start = process.hrtime.bigint();
    let result = 0n;
    for (let i = 0; i < FAST_ITERS; i++) {
        result = fastSum8(1n, 2n, 3n, 4n, 5n, 6n, 7n, 8n);
    }
    let end = process.hrtime.bigint();
    
    const fastTotalNs = Number(end - start);
    const fastTotalMs = fastTotalNs / 1_000_000;
    const fastPerCallNs = fastTotalNs / FAST_ITERS;
    
    console.log(`  Total time:  ${fastTotalMs.toFixed(2).padStart(8)} ms`);
    console.log(`  Per call:    ${fastPerCallNs.toFixed(2).padStart(8)} ns`);

    // === SLOW function ===
    console.log("\nslow_compute (100 calls, 1M iters each):");

    start = process.hrtime.bigint();
    for (let i = 0; i < SLOW_ITERS; i++) {
        result = slowCompute(BigInt(i), COMPUTE_ITERS);
    }
    end = process.hrtime.bigint();
    
    const slowTotalNs = Number(end - start);
    const slowTotalMs = slowTotalNs / 1_000_000;
    const slowPerCallMs = slowTotalMs / SLOW_ITERS;
    
    console.log(`  Total time:  ${slowTotalMs.toFixed(2).padStart(8)} ms`);
    console.log(`  Per call:    ${slowPerCallMs.toFixed(2).padStart(8)} ms`);

    // Prevent optimization
    if (result === -999n) console.log(result);
}

main();
