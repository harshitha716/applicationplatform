FROM debian:stable

# Install required packages
RUN apt-get update

# install awscli
RUN apt-get install awscli -y

# install kubectl
RUN apt-get install -y apt-transport-https ca-certificates curl gnupg
RUN curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.32/deb/Release.key | gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
RUN chmod 644 /etc/apt/keyrings/kubernetes-apt-keyring.gpg
RUN echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.32/deb/ /' | tee /etc/apt/sources.list.d/kubernetes.list
RUN chmod 644 /etc/apt/sources.list.d/kubernetes.list


RUN apt-get update && apt-get install -y kubectl

COPY ./exec.sh /exec.sh

RUN chmod +x /exec.sh

EXPOSE 8099

# CMD 
CMD ["bash", "/exec.sh"]

