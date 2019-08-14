liquidweb_config_path=${PWD}/.lwapi.toml
image=git.liquidweb.com:4567/masre/terraform-provider-liquidweb
dev_image=${image}:dev
mount=-v ${PWD}:/usr/src/terraform-provider-liquidweb
network_name=terraform-provider-liquidweb
network=--network ${network_name}
jaeger_host=-e JAEGER_DISABLED=true -e JAEGER_AGENT_HOST=jaeger -e JAEGER_AGENT_PORT=6831

uid=$(shell id -u)
gid=$(shell id -g)
run_as=--user ${uid}:${gid}

network:
	docker network create ${network_name} || echo "network ${network_name} already exists"

jaeger: network
	docker run -d --rm ${network} --name jaeger -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 jaegertracing/all-in-one:1.13

trace_jaeger:
	#mkfifo pipe
	sh -c "while [ 1 ]; do nc -ul 5775 < pipe | tee outgoing.log | nc -u jaeger 5775 | tee pipe incoming.log; done"

build: clean
	go build

clean:
	rm -f terraform-provider-liquidweb

test:
	go test ./liquidweb

install:
	go install

dev_image:
	docker build --target builder --build-arg uid=${uid} -t ${dev_image} .

shell: dev_image network
	docker run -it --rm --name terraform-provider-liquidweb ${run_as} ${jaeger_host} ${network} ${mount} --entrypoint sh ${dev_image}

init:
	terraform init ${EXAMPLE}

refresh:
	terraform refresh -var liquidweb_config_path=${liquidweb_config_path} -state ${EXAMPLE}/terraform.tfstate -backup=${EXAMPLE}/terraform.tfstate.backup ${EXAMPLE}

plan:
	terraform plan -var liquidweb_config_path=${liquidweb_config_path} -state ${EXAMPLE}/terraform.tfstate ${EXAMPLE}

apply:
  # For proxy use `http_proxy=http://localhost:8080 https_proxy=http://localhost:8080 ...`
	terraform apply -auto-approve -var liquidweb_config_path=${liquidweb_config_path} -state ${EXAMPLE}/terraform.tfstate -backup=${EXAMPLE}/terraform.tfstate.backup ${EXAMPLE} 2>&1 | tee apply.log

destroy:
	terraform destroy -auto-approve -var liquidweb_config_path=${liquidweb_config_path} -state ${EXAMPLE}/terraform.tfstate -backup=${EXAMPLE}/terraform.tfstate.backup ${EXAMPLE}

devplan: build init plan

devapply: build init apply

key:
	ssh-keygen -N '' -C devkey -f ${EXAMPLE}/devkey

image:
	docker build --build-arg uid=${uid} --target builder -t ${image} .

push_image:
	docker push ${image}
