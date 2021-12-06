# GoApiService

一个可以自动构建CURD控制器的go-api服务

## 引入的库

> [gin-gonic/gin](https://github.com/gin-gonic/gin)   【Gin框架】
>
> [go-playground/validator](https://github.com/go-playground/validator)   【validator表单验证器】
>
> [concurrent-map](https://github.com/orcaman/concurrent-map)   【concurrent-map替换原生map解决并发读写】
>
> [gorm.io/gorm](https://gorm.io/gorm)   【Gorm数据查询工具】
>
> [go-redis/redis](https://github.com/go-redis/redis)   【go-redis缓存】
>
> [spf13/viper](https://github.com/spf13/viper)   【配置读取工具】
>
> [uber-go/zap](https://github.com/uber-go/zap)   【zap日志库】
>
>

## 通知邮箱

    osxcoin@protonmail.com

## 打包编译

- 打包期权合约多少秒后消费服务（需要检查redis配置）

```shell
GOOS=linux GOARCH=amd64 go build -o optionContractService service/option_contract/option_contract.go
```

    Redis部分设置
    修改配置文件redis.conf，找到Event notification部分。
    将notify-keyspace-events Ex的注释打开或者添加该配置，其中E代表Keyevent，此种通知会返回key的名字，x代表超时事件。
    如果notify-keyspace-events ""配置没有被注释的话要注释掉，否则不会生效。
    保存后重启redis，一定要使用当前配置文件重启，例如src/redis-server redis.conf

- 打包火币网socket服务

```shell
GOOS=linux GOARCH=amd64 go build -o huobiService service/huobi/socket.go
```

- 打包api接口

```shell
GOOS=linux GOARCH=amd64 go build -o apiService main.go
```

## 交易文档

### 币币交易

### 买入公式

## tips: 当前限价(限价交易取手动输入值，市价交易取K线图值，即火币网获取)

```shell
当前可用余额*百分比/当前限价=买入货币
```

### 卖出公式

## tips: 当前限价(限价交易取手动输入值，市价交易取K线图值，即火币网获取)

```shell
卖出百分比算法：市价（限价）*卖出数量=所得货币
```

## 支付相关

## 回调参数

```json
{
  "merchant_no": "070255",
  "timestamp": "1637661524",
  "sign_type": "MD5",
  "params": {
    "merchant_ref": "I2021112317564133689",
    "system_ref": "TQ1637661402N61AT",
    "amount": "100.00",
    "pay_amount": "100.00",
    "fee": "6.00",
    "status": 1,
    "success_time": 1637661523,
    "extend_params": "",
    "product": "ThaiQR"
  },
  "sign": "2d1ca4805be8b51f62f49f6fad1c879e"
}
```

## 接口文档

### 公共请求字段

- 可以放在header头部

| 字段名      | 类型   | 说明                  |
| ----------- | ------ | --------------------- |
| token | string | 用户登录时获取的token |

### 1、登录模块

#### 1.1、登录接口

- 请求地址

```url
/v1/api/user/login
```

- 请求参数

| 字段名   | 类型   | 说明   |
| -------- | ------ | ------ |
| email | string | 邮箱 |
| password | string | 密码   |

```json
{
  "email": "mail@54zm.com",
  "password": "123456"
}
```

- 返回参数

| 字段名     | 类型     | 说明      |
| ---------- | -------- | --------- |
| code       | string      | 错误代码  |
| msg       | string |     消息提示      |
| reqId       | string |     请求Id     |
| result.token | string   | token钥匙 |
| result.msg   | int   | 用户id  |
| success   | bool   | 成功/失败  |

```json
{
  "code": "200",
  "msg": "success.",
  "reqId": "94ee73c7-be1c-4ed0-a9ed-f0c79a990247",
  "result": {
    "token": "d9433050-816c-4faf-85b3-a8a0fee2006b",
    "uid": 503
  },
  "success": true
}
```

---

