FROM registry.access.redhat.com/ubi9/go-toolset:1.21 AS builder

USER 0

RUN yum -y install yum-utils
RUN yum-config-manager --enable ubi-9-baseos-source

WORKDIR /elfutils-source
RUN yumdownloader --source elfutils
RUN yum -y install cpio
RUN rpm2cpio elfutils-0.191-4.el9.src.rpm | cpio -iv
RUN tar xjvf elfutils-0.191.tar.bz2
WORKDIR /elfutils-source/elfutils-0.191
RUN ./configure --disable-debuginfod
RUN make install

WORKDIR /libbpf-source
RUN yumdownloader --source libbpf
RUN rpm2cpio libbpf-1.4.0-1.el9.src.rpm | cpio -iv
RUN tar xf ./linux-*el9.tar.xz
WORKDIR /libbpf-source/linux-5.14.0-473.el9/tools/lib/bpf
RUN make install_headers
RUN prefix=/usr BUILD_STATIC_ONLY=y make install
WORKDIR /libbpf-source/linux-5.14.0-473.el9/tools/bpf
RUN make bpftool

# rpmautospec requires epel-release
RUN yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm
RUN yum update -y
RUN yum -y install clang rpm-build llvm-devel
#rpmautospec as bug for #1224
# for cpuid on x86, for rpm build
RUN if [ $(uname -i) == "x86_64" ]; then \
    yum install -y https://www.rpmfind.net/linux/centos-stream/9-stream/AppStream/x86_64/os/Packages/http-parser-2.9.4-6.el9.x86_64.rpm && yum install -y cpuid; \
    fi
RUN yum clean all
