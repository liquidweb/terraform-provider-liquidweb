liquidweb_config_path=${PWD}/.lwapi.toml

build:
	go build

ensure:
	dep ensure

install:
	go install

init:
	terraform init

plan:
	terraform plan -var liquidweb_config_path=${storm_config_path}

apply:
	terraform apply -auto-approve -var liquidweb_config_path=${storm_config_path}

destroy:
	terraform destroy -auto-approve -var liquidweb_config_path=${storm_config_path}

devrun: build init apply

key:
	ssh-keygen -N '' -C devkey -f devkey
