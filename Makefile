PROJECT=$(shell basename $(shell pwd))
TAG=ghcr.io/johnjones4/${PROJECT}
VERSION=$(shell date +%s)

LOGHANDLER_TAG=${TAG}-loghandler
LOGHANDLER_DOCKERFILE=./loghandler/Dockerfile

POLLER_SUB_TAG=${TAG}-poller
POLLER_DOCKERFILE=./poller/Dockerfile

SERVER_SUB_TAG=${TAG}-server
SERVER_DOCKERFILE=./server/Dockerfile

.PHONY: loghandler poller server ui

info:
	echo ${PROJECT}

ci: loghandler poller server ui
	
loghandler:
	docker build -t ${LOGHANDLER_TAG} -f ${LOGHANDLER_DOCKERFILE} .
	docker push ${LOGHANDLER_TAG}:latest
	docker image rm ${LOGHANDLER_TAG}:latest

poller:
	docker build -t ${POLLER_SUB_TAG} -f ${POLLER_DOCKERFILE} .
	docker push ${POLLER_SUB_TAG}:latest
	docker image rm ${POLLER_SUB_TAG}:latest

server:
	docker build -t ${SERVER_SUB_TAG} -f ${SERVER_DOCKERFILE} .
	docker push ${SERVER_SUB_TAG}:latest
	docker image rm ${SERVER_SUB_TAG}:latest

ui:
	cd ui && npm install
	cd ui && npm run build
	tar zcvf ui.tar.gz ./ui/build
	git tag ${VERSION}
	git push origin ${VERSION}
	gh release create ${VERSION} ui.tar.gz --generate-notes
