@echo off
setlocal

echo === Building FFI Benchmark ===

where cl >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo ERROR: Run from Visual Studio Developer Command Prompt
    exit /b 1
)

cl /O2 /EHsc bench_native.cpp hotpath.cpp /Fe:bench_native.exe
echo Built bench_native.exe

cl /O2 /LD /EHsc hotpath.cpp /Fe:hotpath.dll
echo Built hotpath.dll

echo.
echo Run: bench_native.exe
echo Run: python bench_python.py (requires SWIG bindings)
echo Run: go run bench_go.go

endlocal
