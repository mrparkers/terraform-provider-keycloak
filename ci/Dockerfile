FROM golang:1.14.4-stretch

RUN apt-get update && \
	apt-get install -y \
		git \
		xvfb \
		apt \
		locales \
		sudo \
		openssh-client \
		ca-certificates \
		tar \
		gzip \
		parallel \
		net-tools \
		netcat \
		unzip \
		zip \
		bzip2 \
		gnupg \
		curl \
		wget && \
		apt-get clean

RUN ln -sf /usr/share/zoneinfo/Etc/UTC /etc/localtime

RUN locale-gen C.UTF-8 || true
ENV LANG=C.UTF-8

ENV TERRAFORM_VERSION "0.12.26"
RUN curl -LO https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
	unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
	sudo mv terraform /usr/local/bin && \
	rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
	terraform version

ENV GOTESTSUM_VERSION "0.4.2"
RUN curl -LO https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_linux_amd64.tar.gz && \
	tar -zxf gotestsum_${GOTESTSUM_VERSION}_linux_amd64.tar.gz && \
	mv gotestsum /usr/local/bin && \
	rm gotestsum_${GOTESTSUM_VERSION}_linux_amd64.tar.gz && \
	gotestsum --version

ENV JQ_VERSION "1.6"
RUN curl -LO https://github.com/stedolan/jq/releases/download/jq-${JQ_VERSION}/jq-linux64 && \
	chmod +x jq-linux64 && \
	mv jq-linux64 /usr/local/bin/jq && \
	jq --version

ENTRYPOINT ["/bin/bash"]
