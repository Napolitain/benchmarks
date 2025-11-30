cmake_minimum_required(VERSION 3.10)
project(hello CXX)

set(CMAKE_CXX_FLAGS_RELEASE "-O3 -flto -march=native -DNDEBUG -s")
set(CMAKE_BUILD_TYPE Release)

add_executable(hello main.cpp)
