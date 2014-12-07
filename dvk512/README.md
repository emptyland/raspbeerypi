## DVK512 介绍

DVK512 是一块 Raspberry Pi Model B+的扩展板,它带有丰富的扩展接口,支持各类 外围模块的接入。

功能：

1. 4个LED
2. 4个按键
3. USB to UART接口
4. RTC实时时钟
5. LCD1602基座

...

## 安装依赖的库

```
sudo yum install libbcm2835 i2c-tools
```

> wiringPi必须要1.18以上版本才支持树莓派B+
> [安装方法](http://wiringpi.com/download-and-install/)

## 开启串口

### 1. 安装cp210x驱动

启动安装在客户端，就是pc端。

驱动下载路径 [Silicon](http://cn.silabs.com/products/mcu/Pages/USBtoUARTBridgeVCPDrivers.aspx)

### 2. 从usb连接到树莓派

*OS X*

1. 确定usb连接是否ok：接口是mini usb与树莓派电源接口micro usb不一样，接口一端连接Mac一端连接DVK512。
2. `ls /dev/tty.SLAB_USBtoUART` 检查终端是否已经建立。
3. `screen /dev/tty.SLAB_USBtoUART 115200` 利用screen连接调试串口。




