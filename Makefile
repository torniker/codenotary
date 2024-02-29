run-app:
	@ go build -o ./build/app . && ./build/app

install-wand:
	@ go build -o ./build/wand ./wand
