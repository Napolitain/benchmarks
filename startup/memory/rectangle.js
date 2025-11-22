#!/usr/bin/env node

const fs = require('fs');
const yaml = require('js-yaml');
const { Command } = require('commander');

function computeRectangleArea(data) {
    const a = data.a || 0;
    const b = data.b || 0;
    const c = data.c || 0;
    const d = data.d || 0;
    
    const width = Math.abs(c - a);
    const height = Math.abs(d - b);
    const area = width * height;
    
    return area;
}

function main() {
    const program = new Command();
    
    program
        .name('rectangle')
        .description('Calculate rectangle area from YAML file')
        .argument('<yaml-file>', 'Path to YAML file containing rectangle coordinates')
        .action((yamlFile) => {
            const start = performance.now();
            
            const fileContents = fs.readFileSync(yamlFile, 'utf8');
            const data = yaml.load(fileContents);
            
            const area = computeRectangleArea(data);
            
            const end = performance.now();
            
            console.log(`Rectangle area: ${area}`);
            console.log(`Time: ${(end - start).toFixed(6)} ms`);
        });
    
    program.parse();
}

main();
