TARGET = hello
SRC = main.rs

all: $(TARGET)

$(TARGET): $(SRC)
	RUSTFLAGS="-C lto=fat -C opt-level=3 -C target-cpu=native -C strip=symbols" rustc -o $(TARGET) $(SRC)

clean:
	rm -f $(TARGET)

.PHONY: all clean
