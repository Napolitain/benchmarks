use std::fs;
use std::time::Instant;

use clap::{Arg, Command};
use serde::Deserialize;

#[derive(Deserialize)]
struct RectangleData {
    a: f64,
    b: f64,
    c: f64,
    d: f64,
}

fn compute_rectangle_area(data: &RectangleData) -> f64 {
    let width = (data.c - data.a).abs();
    let height = (data.d - data.b).abs();
    width * height
}

fn main() {
    let matches = Command::new("rectangle")
        .about("Calculate rectangle area from YAML file")
        .arg(
            Arg::new("yaml-file")
                .help("Path to YAML file containing rectangle coordinates")
                .required(true)
                .index(1),
        )
        .get_matches();

    let yaml_file = matches.get_one::<String>("yaml-file").unwrap();

    let start = Instant::now();

    let file_contents = fs::read_to_string(yaml_file).expect("Error reading file");

    let data: RectangleData = serde_yaml::from_str(&file_contents).expect("Error parsing YAML");

    let area = compute_rectangle_area(&data);

    let elapsed = start.elapsed();

    println!("Rectangle area: {:.2}", area);
    println!("Time: {:.6} ms", elapsed.as_secs_f64() * 1000.0);
}
