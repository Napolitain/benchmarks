TARGET = hello
SRC = main.zig

all: $(TARGET)

$(TARGET): $(SRC)
	zig build-exe -OReleaseFast -fstrip -femit-bin=$(TARGET) $(SRC)

clean:
	rm -f $(TARGET) $(TARGET).o

.PHONY: all clean
