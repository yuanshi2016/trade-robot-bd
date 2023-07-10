#停止 Kubernetes 的服务和组件。运行以下命令停止 kubelet 和其他相关服务
sudo systemctl stop kubelet
sudo systemctl stop docker
#删除 Kubernetes 安装的软件包。运行以下命令来卸载 Kubernetes 相关软件包
sudo apt-get purge kubeadm kubelet kubectl kubernetes-cni
#删除 Kubernetes 的配置文件和数据。运行以下命令来删除相关的目录和文件
sudo rm -rf /etc/kubernetes
sudo rm -rf ~/.kube
#删除 Docker 容器运行时。运行以下命令来删除 Docker 相关软件包和数据
sudo apt-get purge docker-ce docker-ce-cli containerd.io
sudo rm -rf /var/lib/docker
#删除其他相关的依赖项和配置。运行以下命令来删除其他可能存在的依赖项或配置
sudo apt-get autoremove
sudo rm -rf /var/lib/cni/
sudo rm -rf /etc/cni/
sudo rm -rf /var/run/calico
