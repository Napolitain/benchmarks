CC = gcc
CFLAGS = -O3 -flto -march=native -DNDEBUG -s
TARGET = hello
SRC = main.c

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) -o $(TARGET) $(SRC)

clean:
	rm -f $(TARGET)

.PHONY: all clean
