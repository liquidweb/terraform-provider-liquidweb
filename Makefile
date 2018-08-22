build:
	go build

ensure:
	dep ensure

install:
	go install

init:
	terraform init

plan:
	terraform plan -var storm_config_path=${PWD}/.lwapi.toml

apply:
	terraform apply -auto-approve -var "storm_config_path=${PWD}/.lwapi.toml"

destroy:
	terraform destroy -auto-approve

devrun: build init apply

key:
	ssh-keygen -N '' -C devkey -f devkey
