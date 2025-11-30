cmake_minimum_required(VERSION 3.10)
project(hello C)

set(CMAKE_C_FLAGS_RELEASE "-O3 -flto -march=native -DNDEBUG -s")
set(CMAKE_BUILD_TYPE Release)

add_executable(hello main.c)
