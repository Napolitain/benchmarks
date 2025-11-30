TARGET = hello
SRC = main.go

all: $(TARGET)

$(TARGET): $(SRC)
	go build -ldflags="-s -w" -o $(TARGET) $(SRC)

clean:
	rm -f $(TARGET)

.PHONY: all clean
