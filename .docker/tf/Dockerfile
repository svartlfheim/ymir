# FROM hashicorp/terraform:0.14.4
FROM ubuntu:20.04

RUN apt-get update

RUN apt-get install -y ca-certificates curl git

RUN mkdir /usr/local/share/ca-certificates/extra

COPY minica.pem /usr/local/share/ca-certificates/extra/root.cert.crt

RUN update-ca-certificates

COPY --from=hashicorp/terraform:0.14.5 /bin/terraform /bin/terraform
