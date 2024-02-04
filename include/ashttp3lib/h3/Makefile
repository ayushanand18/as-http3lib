BUILD_DIR = /workspaces/as-http3lib/include/deps/quiche/build/debug
INCLUDE_DIR = /workspaces/as-http3lib/include/deps

http3-server:
	g++ -I. -Wall -Werror -pedantic -fsanitize=address -g --std=c++2a \
	-L${BUILD_DIR}/build/quiche-eb2507cd3d43d10d/out/build/ \
	-L${BUILD_DIR}/build/quiche-eb2507cd3d43d10d/out/build/ \
	-L${BUILD_DIR} \
	http3-server.cpp -o http3-server.out -I${INCLUDE_DIR}/quiche/include \
	-I${INCLUDE_DIR}/uthash	${BUILD_DIR}/libquiche.a -lev -ldl -pthread -lm

http3-client:
	g++ -I. -Wall -Werror -pedantic -fsanitize=address -g --std=c++2a \
	-L-L${BUILD_DIR}/build/quiche-eb2507cd3d43d10d/out/build/ \
	-L-L${BUILD_DIR}/build/quiche-eb2507cd3d43d10d/out/build/ \
	-L-L${BUILD_DIR} \
	http3-client.cpp -o http3-client.out -I${INCLUDE_DIR}/quiche/include \
	-I${INCLUDE_DIR}/uthash	${BUILD_DIR}/libquiche.a -lev -ldl -pthread -lm