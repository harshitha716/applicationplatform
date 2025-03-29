FROM oryd/kratos:v1.3.1

WORKDIR /home/ory

COPY ./kratos.yml /home/ory/kratos.yml

COPY ./identity.schema.json /home/ory/identity.schema.json

COPY ./execute.sh /home/ory/kratos.sh

COPY ./mapper.jsonnet /home/ory/mapper.jsonnet
COPY ./okta.jsonnet /home/ory/okta.jsonnet

COPY ./afterregistration.jsonnet /home/ory/afterregistration.jsonnet

ENTRYPOINT [ "sh", "/home/ory/kratos.sh" ]
