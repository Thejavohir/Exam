CURRENT_DIR=$(shell pwd)

proto-gen:
	./script/gen-proto.sh ${CURRENT_DIR}
	ls genproto/*.pb.go | xargs -n1 -IX bash -c "sed -e '/bool/ s/,omitempty//' X > X.tmp && mv X{.tmp.}"

# create_migrate:
# 	migrate create -ext sql -dir migrations -seq create_tablename_table

migrate_up:
	migrate -path migrations/ -database postgres://postgres:Javohir_1@database-1.cxohj5uqpurt.ap-northeast-1.rds.amazonaws.com:5432/customer up

migrate_down:
	migrate -path migrations/ -database postgres://postgres:Javohir_1@database-1.cxohj5uqpurt.ap-northeast-1.rds.amazonaws.com:5432/customer down

pull_submodule:
	git submodule update --init --recursive

update_submodule:
	git submodule update --remote --merge