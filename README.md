# go-pay

go-pay是一个基于grpc微服务的支付网关项目，目前支持微信 JSAPI 支付，后续根据业务需求进行扩展例如：支付宝、银联等。适合在内网环境中部署，提供高效、安全的支付接口。

## 功能特性

- gRPC 通信
- 高性能的远程过程调用，适合微服务架构。
- 定义清晰的 protobuf 接口。

## 快速开始

### 环境依赖

- Go 1.22 或更高版本
- Protobuf 编译器
- 支付相关配置（如商户号、秘钥等）

### 安装依赖

1. 克隆项目到本地：
    ```bash
    git clone https://github.com/supercat0867/go-pay.git
    cd go-pay
    ```
2. 安装 Go 依赖：
    ```bash
    go mod tidy
    ```

### 配置文件

创建 .env 文件并填入以下配置：

```dotenv
DB_USER=root
DB_PASSWORD=123456
DB_HOST=localhost
DB_PORT=3306
DB_NAME=go_pay
```

### 启动服务

项目支持 gRPC 服务 和 HTTP 服务 同时启动：

#### 启动 gRPC 服务

1. gRPC 服务运行在 localhost:9090。
2. 执行以下命令启动服务：
   ```bash
   go run cmd/main.go
   ```

### 客户端示例

#### 获取商户列表：

```go
package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "go-pay/proto"
	"log"
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.NewClient("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// 建立连接
	client := pb.NewMerchantServiceClient(conn)

	resp, err := client.GetMerchants(context.Background(), &pb.GetMerchantsRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Fatalf("调用失败: %v", err)
	}
	log.Printf("调用结果: %s", resp.GetMerchants())
}
```

运行结果：

```bash 
2024/11/15 16:50:12 调用结果: [id:1 name:"测试商户" plant_form:2 plant_form_name:"支付宝" app_id:"123456" mch_id:"123456" cert:"************" cert*******" secret:"************" created_at:"1731652593" id:4 name:"中谷云仓" plant_form:1 plant_form_name:"微信支付" app_id:"wx64d9bf1968f1fdf6" mc71282779" cert:"************" cert_num:"************" secret:"************" created_at:"1731652733" id:7 name:"测试商户" plant_form:1 plant_form_n"微信支付" app_id:"123456" mch_id:"123456" cert:"************" cert_num:"************" secret:"************" created_at:"1731659304"]
```

#### 微信支付JSAPI下单：

```go
package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "go-pay/proto"
	"log"
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.NewClient("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// 建立连接
	client := pb.NewPayServiceClient(conn)

	resp, err := client.GetWechatPrepayInfoJsAPI(context.Background(), &pb.WechatPrepayInfoJsAPIRequest{
		MchId:       "123456", // 商户号
		Amount:      1,
		Description: "测试商品",
		ExpireTime:  "1731674980000",
		Openid:      "ou****************RM",
		OutTradeNo:  "S202411111111111",
		NotifyUrl:   "http://127.0.0.1:8080",
	})
	if err != nil {
		log.Fatalf("调用失败: %v", err)
	}
	log.Printf("调用结果: %s", resp.String())
}

```

运行结果：

```bash 
2024/11/15 16:54:19 调用结果: appId:"wx*******6" timeStamp:"1731660859" nonceStr:"9ef5ef93373efc4802adfe4d5f0bcd90" package:"prepay_id=wx*************00" signType:"RSA" paySign:"Y5sSNS***********MVdrQ=="
```