OS := $(shell uname)

CUR_DIR = /workspaces/as-http3lib/include # replace according to your current directory
BUILD_DIR = $(CUR_DIR)/deps/build
LIB_DIR = $(BUILD_DIR)/debug
INCLUDE_DIR = $(CUR_DIR)/deps/include

INCS = -I$(INCLUDE_DIR)
CFLAGS = -I. -Wall -Werror -pedantic -fsanitize=address -g 

LIBCRYPTO_DIR = $(dir $(shell find ${BUILD_DIR} -name libcrypto.a))
LIBSSL_DIR = $(dir $(shell find ${BUILD_DIR} -name libssl.a))

LDFLAGS = -L$(LIBCRYPTO_DIR) -L$(LIBSSL_DIR) -L$(LIB_DIR)

ifeq ($(OS), Darwin)
# Default prefix of Apple Silicon macOS homebrew.
# Intel will use /usr/local which is included by default.
BREW_INC_DIR = /opt/homebrew/include/
BREW_LIB_DIR = /opt/homebrew/lib/

CFLAGS += -framework Security -I$(BREW_INC_DIR)
LDFLAGS += -L$(BREW_LIB_DIR)
endif

LIBS = $(LIB_DIR)/libquiche.a -lev -ldl -pthread -lm

h1_server:
	@echo "Building Server..."
	g++ --std=c++2a -Wall -Wextra -pedantic h1_server.cpp -o h1_server -lboost_system -pthread -I /workspaces/as-http3lib/include
	@echo "Running Server..."
	./h1_server

client: http3-client.c $(INCLUDE_DIR)/quiche.h $(LIB_DIR)/libquiche.a
	$(CC) $(CFLAGS) $(LDFLAGS) $< -o $@.out $(INCS) $(LIBS)

server: http3-server.cpp $(INCLUDE_DIR)/quiche.h $(LIB_DIR)/libquiche.a
	g++ $(CFLAGS) --std=c++2a $(LDFLAGS) $< -o $@.out $(INCS) $(LIBS)

sample-server: http3-server-sample.cpp $(INCLUDE_DIR)/quiche.h $(LIB_DIR)/libquiche.a
	g++ $(CFLAGS) --std=c++2a $(LDFLAGS) $< -o $@.out $(INCS) $(LIBS)
clean:
	@$(RM) -rf client server http3-client http3-server build/ *.dSYM/
