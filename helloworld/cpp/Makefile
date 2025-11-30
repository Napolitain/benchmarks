CXX = g++
CXXFLAGS = -O3 -flto -march=native -DNDEBUG -s
TARGET = hello
SRC = main.cpp

all: $(TARGET)

$(TARGET): $(SRC)
	$(CXX) $(CXXFLAGS) -o $(TARGET) $(SRC)

clean:
	rm -f $(TARGET)

.PHONY: all clean
