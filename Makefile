liquidweb_config_path=${PWD}/.lwapi.toml

build: clean
	go build

clean:
	rm -f terraform-provider-liquidweb

ensure:
	dep ensure

install:
	go install

init:
	terraform init

plan:
	TF_LOG=trace terraform plan -var liquidweb_config_path=${liquidweb_config_path}

apply:
	terraform apply -auto-approve -var liquidweb_config_path=${liquidweb_config_path}

destroy:
	terraform destroy -auto-approve -var liquidweb_config_path=${liquidweb_config_path}

devplan: build init plan

devapply: build init apply

key:
	ssh-keygen -N '' -C devkey -f devkey
