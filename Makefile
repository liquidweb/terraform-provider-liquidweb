liquidweb_config_path=${PWD}/.lwapi.toml
image=git.liquidweb.com:4567/masre/terraform-provider-liquidweb
dev_image=${image}:dev
mount=-v ${PWD}:/usr/src/terraform-provider-liquidweb

uid=$(shell id -u)
gid=$(shell id -g)
run_as=--user ${uid}:${gid}

build: clean
	go build

clean:
	rm -f terraform-provider-liquidweb

install:
	go install

dev_image:
	docker build --target builder -t ${dev_image} .

shell: dev_image
	docker run -it ${run_as} ${mount} --entrypoint sh ${dev_image}

init:
	terraform init ${EXAMPLE}

refresh:
	terraform refresh -var liquidweb_config_path=${liquidweb_config_path} -state ${EXAMPLE}/terraform.tfstate -backup=${EXAMPLE}/terraform.tfstate.backup ${EXAMPLE}

plan:
	terraform plan -var liquidweb_config_path=${liquidweb_config_path} -state ${EXAMPLE}/terraform.tfstate ${EXAMPLE}

apply:
  # For proxy use `http_proxy=http://localhost:8080 https_proxy=http://localhost:8080 ...`
	terraform apply -auto-approve -var liquidweb_config_path=${liquidweb_config_path} -state ${EXAMPLE}/terraform.tfstate -backup=${EXAMPLE}/terraform.tfstate.backup ${EXAMPLE}

destroy:
	terraform destroy -auto-approve -var liquidweb_config_path=${liquidweb_config_path} -state ${EXAMPLE}/terraform.tfstate -backup=${EXAMPLE}/terraform.tfstate.backup ${EXAMPLE}

devplan: build init plan

devapply: build init apply

key:
	ssh-keygen -N '' -C devkey -f ${EXAMPLE}/devkey

image:
	docker build --target builder -t ${image} .

push_image:
	docker push ${image}
