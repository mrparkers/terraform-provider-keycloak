FROM circleci/golang:1.13.11

ENV TERRAFORM_VERSION "0.12.26"
RUN curl -LO https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
	unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
	sudo mv terraform /usr/local/bin && \
	rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
	terraform version
