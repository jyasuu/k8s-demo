
```sh


sudo docker build -t registry.killer.sh:5000/sun-cipher:latest -t registry.killer.sh:5000/sun-cipher:v1-docker .
sudo docker push registry.killer.sh:5000/sun-cipher:latest
sudo docker push registry.killer.sh:5000/sun-cipher:v1-docker

sudo apt-get -y install podman

podman build -t registry.killer.sh:5000/sun-cipher:v1-podman .
podman push registry.killer.sh:5000/sun-cipher:v1-podman

podman run -d --name sun-cipher registry.killer.sh:5000/sun-cipher:v1-podman
podman ps
podman logs sun-cipher > /opt/course/11/logs
podman ps > /opt/course/11/containers

```