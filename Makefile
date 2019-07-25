liquidweb_config_path=${PWD}/.lwapi.toml
image=git.liquidweb.com:4567/masre/terraform-provider-liquidweb
dev_image=${image}:dev
mount=-v ${PWD}:/usr/src/terraform-provider-liquidweb

build: clean
	go build -o terraform-provider-liquidweb

clean:
	rm -f terraform-provider-liquidweb

install:
	go install

dev_image:
	docker build --target builder -t ${dev_image} .

shell: dev_image
	docker run -it ${mount} --entrypoint sh ${dev_image}

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
	docker build -f Dockerfile -t git.liquidweb.com:4567/masre/terraform-provider-liquidweb .

push_image:
	docker push git.liquidweb.com:4567/masre/terraform-provider-liquidweb
