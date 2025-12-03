#!/usr/bin/env node

import fs from 'fs';
import yaml from 'js-yaml';
import { Command } from 'commander';

interface RectangleData {
    a?: number;
    b?: number;
    c?: number;
    d?: number;
}

function computeRectangleArea(data: RectangleData): number {
    const a = data.a || 0;
    const b = data.b || 0;
    const c = data.c || 0;
    const d = data.d || 0;
    
    const width = Math.abs(c - a);
    const height = Math.abs(d - b);
    const area = width * height;
    
    return area;
}

function main(): void {
    const program = new Command();
    
    program
        .name('rectangle')
        .description('Calculate rectangle area from YAML file')
        .argument('<yaml-file>', 'Path to YAML file containing rectangle coordinates')
        .action((yamlFile: string) => {
            const start = performance.now();
            
            const fileContents = fs.readFileSync(yamlFile, 'utf8');
            const data = yaml.load(fileContents) as RectangleData;
            
            const area = computeRectangleArea(data);
            
            const end = performance.now();
            
            console.log(`Rectangle area: ${area}`);
            console.log(`Time: ${(end - start).toFixed(6)} ms`);
        });
    
    program.parse();
}

main();
