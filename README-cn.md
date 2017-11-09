*阅读本文的其他语言版本：[English](README.md)。*
#  在 IBM 区块链 Hyperledger Fabric V1 上实现常见交易

本项目致力于帮助开发人员从 Hyperledger Fabric V.6 迁移到 V1。本项目将介绍如何在 IBM 区块链上执行传统的数据存储交易。从表面上看，本项目是一个基于 Web 的待办事项列表应用程序，允许执行浏览、读取、编辑、添加和删除 (BREAD) 操作。


这里展示的待办事项列表应用程序，旨在帮助开发人员了解如何调整业务流程所需的常见交易来使用区块链。区块链不等于比特币。可以说比特币是区块链的第一次应用。作为一种分布式账本，区块链的去中心化、共识性和加密等显著特征对许多企业和组织（包括金融、运输、医疗等）有深远的影响。

## 应用程序通信工作流
![待办事项应用程序登录屏幕](todo-list-fabric-client/assets/comm_flow_image.png)

* 用户将在其浏览器中与待办事项客户端应用程序进行交互。
* 当用户执行任何操作时，客户端应用程序调用服务器应用程序 API，已注册的管理员通过该 API 与 Hyperledger 区块链网络进行交互。
* 读取或写入账本称为提案。提案由待办事项应用程序服务器（通过 SDK）构建，然后发送到一个区块链对等节点。
* 该对等节点将与它的待办事项链代码容器进行通信。链代码将运行/模拟该交易。如果没有问题，它会对该交易进行背书，并将其发回我们的待办事项应用程序。
* 然后，待办事项应用程序（通过 SDK）将背书后的提案发送到订购服务。订购方将来自整个网络的许多提案打包到一个区块中。然后，它将新的区块广播到网络中的对等节点。
* 最后，对等节点会验证该区块并将它写入自己的账本中。该交易现在已经生效，所有后续读取都会反映此更改。

## 前提条件

* [Go](https://golang.org/) - 最新版本
* [Docker](https://www.docker.com/products/overview) - V1.13 或更高版本
* [Docker Compose](https://docs.docker.com/compose/overview/) - V1.8 或更高版本
* [Node.js & npm](https://nodejs.org/en/download/) - node v6.2.0 - V6.10.0（不支持 V7 及更高版本）；您的 Node 安装中包含 npm。
* [xcode](https://developer.apple.com/xcode/) - 仅 OS X 用户需要
* [nvm](https://github.com/creationix/nvm/blob/master/README.markdown) - 如果您想使用 nvm install 命令检索 Node 版本


## 步骤

1. [下载 Docker 映像并获取 hyperledger fabric V1 node sdk 的代码](#1-download-the-docker-images-and-get-the-code-for-hyperledger-fabric-v1-node-sdk)

2. [编辑配置](#2-edit-the-configuration)

3. [启动您的网络](#3-start-your-network)

4. [使用 Node SDK](#4-use-the-node-sdk)

5. [运行待办事项列表 fabric 服务器](#5-run-the-todo-list-fabric-server)

6. [运行待办事项列表 fabric 客户端](#6-run-the-todo-list-fabric-client)

7. [使用待办事项列表应用程序](#7-run-the-todo-list-application)


# 1.下载 Docker 映像并获取 hyperledger fabric V1 node sdk 的代码

`download-dockerimages.sh` 包含用于下载所需的 Docker 映像的代码，设置运行 Hyperledger Fabric V1 的网络需要这些映像。

从工作区中，让该 shell 脚本变得可执行：

```bash
chmod +x download-dockerimages.sh
```

现在运行该脚本。在执行此脚本之前，确保 Docker 正在运行。此过程会花费几分钟的时间，所以请耐心等待：

```bash
./download-dockerimages.sh
```

该脚本执行完成后，您会在终端中看到以下信息：

```bash
===> List out hyperledger docker images
hyperledger/fabric-ca          latest               35311d8617b4        3 weeks ago         240 MB
hyperledger/fabric-ca          x86_64-1.0.0-alpha   35311d8617b4        3 weeks ago         240 MB
hyperledger/fabric-couchdb     latest               f3ce31e25872        3 weeks ago         1.51 GB
hyperledger/fabric-couchdb     x86_64-1.0.0-alpha   f3ce31e25872        3 weeks ago         1.51 GB
hyperledger/fabric-kafka       latest               589dad0b93fc        3 weeks ago         1.3 GB
hyperledger/fabric-kafka       x86_64-1.0.0-alpha   589dad0b93fc        3 weeks ago         1.3 GB
hyperledger/fabric-zookeeper   latest               9a51f5be29c1        3 weeks ago         1.31 GB
hyperledger/fabric-zookeeper   x86_64-1.0.0-alpha   9a51f5be29c1        3 weeks ago         1.31 GB
hyperledger/fabric-orderer     latest               5685fd77ab7c        3 weeks ago         182 MB
hyperledger/fabric-orderer     x86_64-1.0.0-alpha   5685fd77ab7c        3 weeks ago         182 MB
hyperledger/fabric-peer        latest               784c5d41ac1d        3 weeks ago         184 MB
hyperledger/fabric-peer        x86_64-1.0.0-alpha   784c5d41ac1d        3 weeks ago         184 MB
hyperledger/fabric-javaenv     latest               a08f85d8f0a9        3 weeks ago         1.42 GB
hyperledger/fabric-javaenv     x86_64-1.0.0-alpha   a08f85d8f0a9        3 weeks ago         1.42 GB
hyperledger/fabric-ccenv       latest               91792014b61f        3 weeks ago         1.29 GB
hyperledger/fabric-ccenv       x86_64-1.0.0-alpha   91792014b61f        3 weeks ago         1.29 GB
```

克隆 fabric node sdk 的存储库：
```bash
git clone https://github.com/hyperledger/fabric-sdk-node.git
```

首先，签出 `fabric-sdk-node` 存储库的 alpha 分支：
```bash
cd fabric-sdk-node
git checkout v1.0.0-alpha
```

确保您在正确的分支上：
```bash
git branch
```

您会看到以下结果：
```bash
Ishans-MacBook-Pro:fabric-sdk-node ishan$ git branch
* (HEAD detached at v1.0.0-alpha)
  master
```

现在跳回您的工作区目录：
```bash
cd ..
```

在您的工作区中，将 docker-compose-networksetup.yaml 迁移到 fabric-sdk-node 目录中的 test/fixtures 文件夹：

```bash
mv docker-compose-networksetup.yaml fabric-sdk-node/test/fixtures
```

仍在您的工作区中，清空 fabric-sdk-node 目录中的示例链代码源代码：

```bash
rm -rf fabric-sdk-node/test/fixtures/src/github.com/example_cc/*
```

现在将待办事项列表链代码复制到相同文件夹：
```bash
cp todo-list-fabric-server/chaincode/* fabric-sdk-node/test/fixtures/src/github.com/example_cc/
```
> **备注：**如果您想在 hyperledger fabric V1 上运行自己的代码，只需将链代码的代码复制到 fabric-sdk-node/test/fixtures/src/github.com/example_cc 目录中。

# 2.编辑配置

更新 fabric-sdk-node 目录中的 `config.json` 和 `instantiate-chaincode.js` 文件：

```bash
cd fabric-sdk-node/test/integration/e2e
```

使用编辑器打开 `config.json`，并将所有 `grpcs` 实例替换为 `grpc` 实例。

使用编辑器打开 `instantiate-chaincode.js`，并将第 147 行替换为：
```bash
args: ['init'],
```

# 3.启动您的网络

`docker-compose-networksetup.yaml` 包含用于设置网络的配置。

导航到 fabric-sdk-node 目录中的 test/fixtures 文件夹，并运行 docker-compose 文件：

```bash
cd fabric-sdk-node/test/fixtures
docker-compose -f docker-compose-networksetup.yaml up -d
```

完成上述操作后，发出一个 docker ps 命令来查看目前运行的容器。您会看到以下结果：
```bash
CONTAINER ID        IMAGE                        COMMAND                  CREATED             STATUS                       PORTS                                            NAMES
e61cf829f171        hyperledger/fabric-peer      "peer node start -..."3 minutes ago       Up 2 minutes           0.0.0.0:7056->7051/tcp, 0.0.0.0:7058->7053/tcp   peer1
0cc1f5ac24da        hyperledger/fabric-peer      "peer node start -..."3 minutes ago       Up 2 minutes        0.0.0.0:8056->7051/tcp, 0.0.0.0:8058->7053/tcp   peer3
7ab3106e5076        hyperledger/fabric-peer      "peer node start -..."3 minutes ago       Up 3 minutes        0.0.0.0:7051->7051/tcp, 0.0.0.0:7053->7053/tcp   peer0
2bc5c6606e6c        hyperledger/fabric-peer      "peer node start -..."3 minutes ago       Up 3 minutes        0.0.0.0:8051->7051/tcp, 0.0.0.0:8053->7053/tcp   peer2
513be1b46467        hyperledger/fabric-ca        "sh -c 'fabric-ca-..."3 minutes ago       Up 3 minutes        0.0.0.0:8054->7054/tcp                           ca_peerOrg2
741c363ba34a        hyperledger/fabric-orderer   "orderer"                3 minutes ago       Up 3 minutes        0.0.0.0:7050->7050/tcp                           orderer0
abaae883eb13        couchdb                      "tini -- /docker-e..."3 minutes ago       Up 3 minutes        0.0.0.0:5984->5984/tcp                           couchdb
2c2d51fe88c0        hyperledger/fabric-ca        "sh -c 'fabric-ca-..."3 minutes ago       Up 3 minutes        0.0.0.0:7054->7054/tcp                           ca_peerOrg1

```

# 4.使用 Node SDK

返回到 `fabric-sdk-node` 目录的根目录，并将 grpc 依赖项 `"grpc": "1.1.2"` 添加到 package.json

将 node 模块安装到您的 SDK 存储库中。
```bash
npm install
npm install -g gulp
# if you get a "permission denied" error, then try with sudo
sudo npm install -g gulp
```
最后，构建 fabric-ca 客户端：
```bash
gulp ca
```

删除在前面的运行中可能缓存的键值存储和 hfc 工件：
```bash
rm -rf /tmp/hfc-*
rm -rf ~/.hfc-key-store
```

将 `test/unit` 中的 `util.js` 中的链代码版本从 V0 更新为 V1
```bash
module.exports.END2END = {
	channel: 'mychannel',
	chaincodeId: 'end2end',
	chaincodeVersion: 'v1'
};
```

### 创建通道

Hyperledger Fabric 通道是两个或多个特定网络成员之间的通信的私有“子集”，旨在执行私有且机密的交易。通道由成员（组织）、每个成员的锚对等节点、共享账本、链代码应用程序和订购服务节点来定义。网络上的每个交易都在一个通道上执行，每方必须经过验证和授权才能在该通道上执行交易。成员服务提供者 (MSP) 会为每个加入通道的对等节点分配它们自己的身份，向通道对等节点和服务验证每个对等节点。

现在，利用 SDK 测试程序创建一个名为 `mychannel` 的通道。从 `fabric-sdk-node` 目录运行以下代码：
```bash
node test/integration/e2e/create-channel.js
```

### 加入通道
将创始区块 `mychannel.block` 传递到订购服务，并将对等节点加入您的通道中：
```bash
node test/integration/e2e/join-channel.js
```

### 安装链代码
将待办事项列表源代码安装在对等节点的文件系统上：
```bash
node test/integration/e2e/install-chaincode.js
```

### 实例化链代码
创建代办事项列表容器：
```bash
node test/integration/e2e/instantiate-chaincode.js
```

# 5.运行待办事项列表 fabric 服务器

导航到 `todo-list-fabric-server` 目录的根目录。

将 node 模块安装到您的 fabric 服务器存储库中。
```bash
npm install
```

运行该服务器：
```bash
node server.js
```

向 `/enrollAdmin` 端点发出一个 get 请求，以在链代码上注册管理员：

您会看到以下响应：
```
{
  "message": "Admin Enrolled! "
}
```
# 6.运行待办事项列表 fabric 客户端

在新的终端中，导航到 `todo-list-fabric-client/web` 目录的根目录。

为了让基于 Web 的待办事项列表应用程序正常工作，必须从 Web 服务器运行它。为了让该应用程序正常运行，不需要公开提供这个服务器。

在 Mac 上，一种常见方法是使用内置的 PHP 安装来就地运行 Web 服务器。
运行 PHP Web 服务器：
```bash
php -S localhost:8081
```

在 Windows 上，可以使用 [XAMPP](https://www.apachefriends.org/index.html)
# 7.使用待办事项列表应用程序

使用链接 `http://localhost:8081` 将该 Web 应用程序加载到浏览器中。您会看到一个登录屏幕。登录对话框包含一个 IBM 徽标。按住 Alt 并单击该徽标，以便将数据预先加载到区块链中。唯一表明此操作已完成的指标是开发人员控制台中的交易 ID。

> 尽管不是很冗长，但是每次在区块链本身上执行更改，都会在浏览器的开发人员控制台中显示来自 IBM 区块链的交易 ID。在使用待办事项列表应用程序时，打开开发人员控制台可能很有用。

![待办事项登录屏幕](todo-list-fabric-client/assets/todo-authentication.png)

在默认数据中创建了 3 个帐户。这些帐户以“用户名:密码”形式表示为：

- krhoyt:abc123
- abtin:abc123
- peter:abc123

可以使用任何这些帐户进行登录，以便浏览、读取、编辑、添加和删除待办事项。

![待办事项列表](todo-list-fabric-client//assets/todo-list.png)

- 要创建待办事项列表项，请单击标为 "+" 的红色按钮。将鼠标悬停在此按钮上，就会显示用于创建“位置”的额外按钮。
- 要编辑待办事项列表项，请单击您打算编辑的项，并修改相应字段，使之与您想要的值匹配。没有 Save 按钮，因为所有更改会立刻提交到区块链。
- 要删除待办事项列表项，请将鼠标移到任何项上，然后单击垃圾桶图标。
- 要将待办事项列表项转发给另一个人，请将鼠标移到任何项上，然后单击出现的箭头图标。系统中的其他用户的列表将会出现。选择一个名称。
- 要注销应用程序，请单击一个方框内含一个箭头的图标。该图标位于屏幕的右上角。
- 使用上面的帐户信息，再次使用不同的帐户登录到该应用程序，以查看转发到系统中的其他用户的待办事项。

# 附加资源
以下是一个附加区块链资源列表：
* [IBM 区块链基础](https://www.ibm.com/blockchain/what-is-blockchain.html)
* [Hyperledger Fabric 文档](http://fabric-rtd.readthedocs.io/en/latest/getting_started.html)
* [GitHub 上的 Hyperledger Fabric 代码](https://github.com/hyperledger/fabric)
* [Hyperledger Fabric Composer](https://hyperledger.github.io/composer/)
* [如何迁移基于 Fabric V0.6 的链代码以在最新的 Fabric V1.0 上运行](https://developer.ibm.com/blockchain/2017/03/17/migrate-fabric-v0-6-based-chaincode-run-latest-fabric-v1-0/)

# 故障排除

* 如果在运行 docker-compose 时看到消息 `containerID already exists`，则需要删除现有的容器。此命令将删除所有容器；“而不是”您的映像：
```bash
docker rm -f $(docker ps -aq)
```

* 在运行 `create-channel.js` 时，如果看到错误 `private key not found`，可以尝试清除已缓存的键值存储：
```bash
rm -rf /tmp/hfc-*
rm -rf ~/.hfc-key-store
```

* 浏览器中的开发人员控制台，是排除在运行客户端应用程序时可能出现的任何问题的关键。查找错误的第一个地方是 /web/script/blockchain.js 文件中的链代码 ID 和 URL 的值。

# 参考资料
* 本示例基于使用 Hyperledger Fabric V0.6 的待办事项列表应用程序 [Hyperledger Fabric V0.6](https://github.com/IBM/todo-list-fabric)。

# 许可
[Apache 2.0](LICENSE)
