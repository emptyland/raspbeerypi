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

## 使用RTC时钟

树莓派本身不带时钟电池，导致每次重启时间都归0（1970年1月1日0时0分0秒）。DVK512扩展板提供了一个RTC时钟源与CR1220(直径12毫米)纽扣电池座。

查找I2C地址：

`sudo i2cdetect -y 1`

输出结果类似：

```
     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
00:          -- -- -- -- -- -- -- -- -- -- -- -- --
10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
50: -- 51 -- -- -- -- -- -- -- -- -- -- -- -- -- --
60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
70: -- -- -- -- -- -- -- --
```

注册PF8563:

```
sudo sh -c 'echo pcf8563 0x51 > /sys/class/i2c-adapter/i2c-1/new_device'
```

* 读I2C硬件时钟：`hwclock –r`
* 将系统时钟写入PCF863 `hwclock –w`
* 设置系统时钟与RTC同步：`hwclock –s`







