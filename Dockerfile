FROM centos:7.6.1810

COPY taichi /usr/local/bin/taichi
EXPOSE 3000
ENTRYPOINT ["/usr/local/bin/taichi"]
