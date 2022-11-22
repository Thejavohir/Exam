CURRENT_DIR=$(shell pwd)

proto-gen:
	./script/gen-proto.sh ${CURRENT_DIR}
	ls genproto/*.pb.go | xargs -n1 -IX bash -c "sed -e '/bool/ s/,omitempty//' X > X.tmp && mv X{.tmp.}"

migrate_up:
	migrate -path migrations/ -database postgres://postgres:1234@localhost:5432/customer up

migrate_down:
	migrate -path migrations/ -database postgres://postgres:1234@localhost:5432/customer down

pull_submodule:
	git submodule update --init --recursive

update_submodule:
	git submodule update --remote --merge
