#更新系统软件包
sudo apt update
sudo apt upgrade -y
# 下载 Google Cloud 公开签名秘钥：
sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
# 添加 Kubernetes apt 仓库
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list

#安装 Docker 容器运行时。运行以下命令来安装 Docker 启动并启用 Docker 服务
sudo apt install docker.io -y
sudo systemctl start docker
sudo systemctl enable docker
#添加 Kubernetes 的软件仓库。运行以下命令来添加 Kubernetes 官方软件仓库
sudo apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"
sudo apt-add-repository "deb https://mirrors.tuna.tsinghua.edu.cn/google-cloud-packages/apt cloud-sdk main"
#安装 Kubernetes 组件。运行以下命令来安装 kubelet、kubeadm 和 kubectl
sudo apt install kubelet kubeadm kubectl -y
#预先拉取镜像
kubeadm config images pull --config kubeadm-config.yaml

#初始化 Kubernetes 控制平面。运行以下命令来初始化 Kubernetes 控制平面组件
sudo kubeadm init

#完成初始化后，按照命令输出的指示，将提供的 kubectl 配置文件复制到您的用户目录下.
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

#安装网络插件。运行以下命令来安装网络插件，例如 Calico
kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml

#等待一段时间，直到所有的 Kubernetes 组件都处于运行状态。您可以使用以下命令检查组件的状态
kubectl get pods --all-namespaces


