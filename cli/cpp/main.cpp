#include <argparse/argparse.hpp>
#include <ryml.hpp>
#include <ryml_std.hpp>
#include <c4/format.hpp>

#include <chrono>
#include <cmath>
#include <fstream>
#include <iostream>
#include <sstream>
#include <string>

struct RectangleData {
    double a, b, c, d;
};

double computeRectangleArea(const RectangleData& data) {
    double width = std::abs(data.c - data.a);
    double height = std::abs(data.d - data.b);
    return width * height;
}

int main(int argc, char* argv[]) {
    argparse::ArgumentParser program("rectangle");
    program.add_argument("yaml_file")
        .help("YAML file containing rectangle coordinates");

    try {
        program.parse_args(argc, argv);
    } catch (const std::exception& err) {
        std::cerr << err.what() << std::endl;
        std::cerr << program;
        return 1;
    }

    auto start = std::chrono::high_resolution_clock::now();

    std::string yaml_file = program.get<std::string>("yaml_file");

    std::ifstream file(yaml_file);
    if (!file.is_open()) {
        std::cerr << "Error reading file: " << yaml_file << std::endl;
        return 1;
    }

    std::stringstream buffer;
    buffer << file.rdbuf();
    std::string contents = buffer.str();

    ryml::Tree tree = ryml::parse_in_arena(ryml::to_csubstr(contents));
    ryml::ConstNodeRef root = tree.rootref();

    RectangleData data;
    root["a"] >> data.a;
    root["b"] >> data.b;
    root["c"] >> data.c;
    root["d"] >> data.d;

    double area = computeRectangleArea(data);

    auto end = std::chrono::high_resolution_clock::now();
    auto elapsed = std::chrono::duration<double, std::milli>(end - start);

    std::cout << "Rectangle area: " << std::fixed << area << std::endl;
    std::cout << "Time: " << elapsed.count() << " ms" << std::endl;

    return 0;
}
