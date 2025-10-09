CC = gcc
CFLAGS = -Wall -Wextra -std=c11 -I../../include
LDFLAGS = -L../../native/$(shell uname -s | tr '[:upper:]' '[:lower:]')_$(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -Wl,-rpath=../../native/$(shell uname -s | tr '[:upper:]' '[:lower:]')_$(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
LIBS = -lextractous_ffi

$(info [debug] LDFLAGS = $(LDFLAGS))
$(info [debug] LIBS    = $(LIBS))
$(info [debug] CFLAGS  = $(CFLAGS))

TARGET = smoke
SRC = smoke.c

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) $(LDFLAGS) -o $@ $< $(LIBS)

run: $(TARGET)
	LD_LIBRARY_PATH=.:$(LD_LIBRARY_PATH) ../../native/$(TARGET)

clean:
	rm -f $(TARGET)
	rm -rf test_files

.PHONY: all clean run
rebuild: clean all
