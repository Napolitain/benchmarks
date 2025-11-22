#!/usr/bin/env python3
import argparse
import yaml
import time

def compute_rectangle_area(data):
    a = data.get('a', 0)
    b = data.get('b', 0)
    c = data.get('c', 0)
    d = data.get('d', 0)
    
    width = abs(c - a)
    height = abs(d - b)
    area = width * height
    
    return area

def main():
    parser = argparse.ArgumentParser(description='Calculate rectangle area from YAML file')
    parser.add_argument('yaml_file', help='Path to YAML file containing rectangle coordinates')
    args = parser.parse_args()
    
    start = time.perf_counter()
    
    with open(args.yaml_file, 'r') as f:
        data = yaml.safe_load(f)
    
    area = compute_rectangle_area(data)
    
    end = time.perf_counter()
    
    print(f"Rectangle area: {area}")
    print(f"Time: {(end - start) * 1000:.6f} ms")

if __name__ == "__main__":
    main()
