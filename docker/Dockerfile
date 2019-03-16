FROM hashicorp/terraform:0.11.13

RUN apk add --update make bash go build-base qemu-system-x86_64 libvirt libvirt-dev libvirt-daemon pkgconfig

WORKDIR $GOPATH/src/github.com/dmacvicar/terraform-provider-libvirt

ARG GIT_REPOSITORY
ARG GIT_BRANCH

ENV GIT_REPOSITORY=${GIT_REPOSITORY:-https://github.com/dmacvicar/terraform-provider-libvirt}
ENV GIT_BRANCH=${GIT_BRANCH:-master}

RUN git clone $GIT_REPOSITORY ./ && \
  git checkout $GIT_BRANCH && \
  go get -u

ENTRYPOINT ["docker/testacc.sh"]
