TARGET = hello
SRC = Main.java
CLASS = Main.class

all: $(CLASS)

$(CLASS): $(SRC)
	javac $(SRC)

run: $(CLASS)
	java Main

clean:
	rm -f $(CLASS)

.PHONY: all run clean
