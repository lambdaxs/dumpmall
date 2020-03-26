all:build move restart

dev:build move_mac restart

move_mac:
        mv dump_server /Users/xiaos/servers/

build:main.go
        go build -o dump_server

move:
        mv dump_server /root/servers/

restart:
        supervisorctl restart dump_server