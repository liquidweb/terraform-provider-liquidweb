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
	terraform init ${EXAMPLE}

refresh:
	terraform refresh -var liquidweb_config_path=${liquidweb_config_path} -state ${EXAMPLE}/terraform.tfstate -backup=${EXAMPLE}/terraform.tfstate.backup ${EXAMPLE}

plan:
	terraform plan -var liquidweb_config_path=${liquidweb_config_path} ${EXAMPLE}

apply:
	terraform apply -auto-approve -var liquidweb_config_path=${liquidweb_config_path} -state ${EXAMPLE}/terraform.tfstate -backup=${EXAMPLE}/terraform.tfstate.backup ${EXAMPLE}

destroy:
	terraform destroy -auto-approve -var liquidweb_config_path=${liquidweb_config_path} -state ${EXAMPLE}/terraform.tfstate -backup=${EXAMPLE}/terraform.tfstate.backup ${EXAMPLE}

devplan: build init plan

devapply: build init apply

key:
	ssh-keygen -N '' -C devkey -f ${EXAMPLE}/devkey

image:
	docker build -f terraform.Dockerfile -t git.liquidweb.com:4567/masre/terraform-provider-liquidweb .

push_image:
	docker push git.liquidweb.com:4567/masre/terraform-provider-liquidweb
