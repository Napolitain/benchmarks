# Startup Benchmarks

Two benchmarks implemented in Python, JavaScript (Node.js), and Go:

1. **Compute (Bubblesort)**: Sorts 10 integers using bubble sort algorithm
2. **Memory (Rectangle Area)**: Parses a YAML file with rectangle coordinates (a, b, c, d) and computes the area

## Structure

```
compute/        - CPU-intensive bubblesort benchmark
memory/         - I/O and memory benchmark (YAML parsing + computation)
```

## Setup

### Python
```bash
pip install pyyaml
```

### JavaScript/Node.js
```bash
npm install
```
Dependencies: `js-yaml`, `commander`

### Go
```bash
go mod download
```
Dependencies: `cobra`, `goccy/go-yaml`

## Running the Benchmarks

### Compute - Bubblesort

**Python:**
```bash
python3 compute/bubblesort.py
```

**JavaScript:**
```bash
node compute/bubblesort.js
```

**Go:**
```bash
go run compute/bubblesort.go
```

### Memory - Rectangle Area Calculator

**Python (argparse + pyyaml):**
```bash
python3 memory/rectangle.py memory/test_rectangle.yaml
```

**JavaScript (commander + js-yaml):**
```bash
node memory/rectangle.js memory/test_rectangle.yaml
```

**Go (cobra + goccy/go-yaml):**
```bash
go run memory/rectangle.go memory/test_rectangle.yaml
```

## Test YAML Format

The `memory/test_rectangle.yaml` file contains:
```yaml
a: 0
b: 0
c: 10
d: 5
```

Where `a` and `b` are coordinates of one corner, and `c` and `d` are coordinates of the opposite corner.
The area is calculated as: `|c - a| * |d - b|`
