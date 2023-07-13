### 申明: 个人作品，不可用于商业交易
* iFortune一款网格量化交易服务，提供网格策略交易、趋势网格、反向网格等免费量化交易策略，已支持接入Huobi、OKex、Binance等主流交易平台，覆盖USDT、BTC、ETH等多种交易对，通过API一键配置策略，智能自动交易。同时设置专业完善的风控体系保障账户安全，保护策略安全，最大程度降低量化投资风险。
#### 采用的技术
* 基于go-kratos微服务实战
* 基于k8s编排，Jenkins持续交付
* [配套基于Flutter开发Android、iOS、Web](https://github.com/RonadoLong/wq-fotune.git)
* [APP体验版本:https://www.pgyer.com/gnTX](https://www.pgyer.com/gnTX)
* [web体验版本:https://yun.local100.com/test/web/#/](https://yun.local100.com/test/web/#/)
#### 简易架构图


#### 前端部分页面展示
![](./resource/images/WechatIMG33.jpeg)
![](./resource/images/WechatIMG34.jpeg)
![](./resource/images/WechatIMG36.jpeg)
![](./resource/images/WechatIMG37.jpeg)
![](./resource/images/WechatIMG38.jpeg)
![](./resource/images/WechatIMG39.jpeg)

#### 当前部署使用说明
* 目前是本地服务器自建k8s，一个master，两个worker。基于frpc进行内网穿透，流量转发；
* 使用rancher管理k8s集群；
* 使用Jenkins pipeline进行持续化构建；

## 环境安装


# docker相关
```shell
#显示容器信息
docker ps --format "table {{.ID}}\t{{.Names}}\t{{.Status}}\t{{.Ports}}"

#镜像推送
docker tag quote-svc:quote-svc harbor.local100.com/mateforce/quote-svc:quote-svc
docker push harbor.local100.com/mateforce/trade_rebot_builder:latest
```
#### 安装rancher
```shell
# system-default-registry地址
registry.cn-hangzhou.aliyuncs.com
# 关闭swap
swapoff -a
# 关闭防火墙
systemctl stop firewalld && systemctl disable firewalld
# 完全清理脚本 - 仅在重复安装出问题后
curl -LO https://github.com/rancherlabs/support-tools/raw/master/extended-rancher-2-cleanup/extended-cleanup-rancher2.sh
bash extended-cleanup-rancher2.sh

docker run -d --privileged --restart=unless-stopped -p 2080:80 -p 2043:443 -v /opt/rancher:/var/lib/rancher rancher/rancher:v2.7.5
密码:RNntiyObLh8WB62Q
```
#### 安装Harbor
```shell
# 下载安装包
wget https://github.com/goharbor/harbor/releases/download/v2.2.0/harbor-offline-installer-v2.2.0.tgz

# 解压
tar -zxvf harbor-offline-installer-v2.2.0.tgz
mv harbor /var/local/harbor
cd /var/local/harbor

# 修改配置文件
cp harbor.yml.tmpl harbor.yml
`
hostname: 0.0.0.0
http:
  port: 8086
harbor_admin_password: admin
database:
  password: admin
data_volume: /harbor_data
`
# 执行安装

cd /var/local/harbor/
./install.sh

# docker-compose启动
docker-compose up -d
# docker-compose重启
docker-compose down -v


# 配置加速
echo > /etc/docker/daemon.json
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "insecure-registries": ["harbor.local100.com:8086","harbor.local100.com","harbor.local100.com:8086","10.10.1.100:8086","0.0.0.0"],
  "registry-mirrors": [
        "https://mirrors.sjtug.sjtu.edu.cn",
        "https://mirror.ccs.tencentyun.com",
        "https://docker.mirrors.ustc.edu.cn",
        "https://hub-mirror.c.163.com"
    ]
}
EOF
systemctl daemon-reload && systemctl restart docker && systemctl restart harbor

docker login -u admin -p admin harbor.local100.com
```
#### jenkins运行
```shell
docker run -u root -d -p 18080:8080 -p 50000:50000 -v /var/jenkins:/var/jenkins_home -v /var/run/docker.sock:/var/run/docker.sock jenkins/jenkins
```