

> 假设使用Pidora

## 脚针定义

> 由`gpio readall`输出

<pre>
 |     |     |    3.3v |      |   |  1 || 2  |   |      | 5v      |     |     |
 |   2 |   8 |   SDA.1 | ALT0 | 1 |  3 || 4  |   |      | 5V      |     |     |
 |   3 |   9 |   SCL.1 | ALT0 | 1 |  5 || 6  |   |      | 0v      |     |     |
 |   4 |   7 | GPIO. 7 |   IN | 1 |  7 || 8  | 1 | ALT0 | TxD     | 15  | 14  |
 |     |     |      0v |      |   |  9 || 10 | 1 | ALT0 | RxD     | 16  | 15  |
 |  17 |   0 | GPIO. 0 |   IN | 0 | 11 || 12 | 0 | IN   | GPIO. 1 | 1   | 18  |
 |  27 |   2 | GPIO. 2 |   IN | 0 | 13 || 14 |   |      | 0v      |     |     |
 |  22 |   3 | GPIO. 3 |   IN | 0 | 15 || 16 | 0 | IN   | GPIO. 4 | 4   | 23  |
 |     |     |    3.3v |      |   | 17 || 18 | 0 | IN   | GPIO. 5 | 5   | 24  |
 |  10 |  12 |    MOSI | ALT0 | 0 | 19 || 20 |   |      | 0v      |     |     |
 |   9 |  13 |    MISO | ALT0 | 0 | 21 || 22 | 0 | IN   | GPIO. 6 | 6   | 25  |
 |  11 |  14 |    SCLK | ALT0 | 0 | 23 || 24 | 1 | ALT0 | CE0     | 10  | 8   |
 |     |     |      0v |      |   | 25 || 26 | 1 | ALT0 | CE1     | 11  | 7   |
 |   0 |  30 |   SDA.0 | ALT0 | 1 | 27 || 28 | 1 | ALT0 | SCL.0   | 31  | 1   |
 |   5 |  21 | GPIO.21 |   IN | 1 | 29 || 30 |   |      | 0v      |     |     |
 |   6 |  22 | GPIO.22 |   IN | 1 | 31 || 32 | 0 | IN   | GPIO.26 | 26  | 12  |
 |  13 |  23 | GPIO.23 |   IN | 0 | 33 || 34 |   |      | 0v      |     |     |
 |  19 |  24 | GPIO.24 |   IN | 0 | 35 || 36 | 0 | IN   | GPIO.27 | 27  | 16  |
 |  26 |  25 | GPIO.25 |   IN | 0 | 37 || 38 | 0 | IN   | GPIO.28 | 28  | 20  |
 |     |     |      0v |      |   | 39 || 40 | 0 | IN   | GPIO.29 | 29  | 21  |
 +-----+-----+---------+------+---+----++----+---+------+---------+-----+-----+
 | BCM | wPi |   Name  | Mode | V | Physical | V | Mode | Name    | wPi | BCM |
 +-----+-----+---------+------+---+--B Plus--+---+------+---------+-----+-----+
</pre>


## 基本设置

### sudo免密码

1. 确定当前用户是否在`wheel`组中

	```
	$ groups
	```
2. 为`wheel`组设定免密码

	```
	$ sudo sh -c 'echo ""%wheel ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers'
	```
> 由于给所有`wheel`组的成员设置了免密码，可能不太安全。

### 配置vim

```
curl 'https://raw.githubusercontent.com/emptyland/scripts/master/conf/comm/vimrc' > ~/.vimrc
```


## 如何安装fish shell

fish shell在pidroa的官方源里没有，需要从源代码安装。

### 安装 g++
```
sudo yum install gcc-c++
```

### 安装fish依赖的curses库
```
sudo yum install ncurses-devel
```

### 获取fish代码
```
cd ~
wget 'http://fishshell.com/files/2.1.1/fish-2.1.1.tar.gz'
```

### 变异并安装
```
tar vxzf fish-2.1.1.tar.gz
cd fish-2.1.1
./configure --prefix=/usr
make
sudo make install
```

> 由于pi的性能限制，编译时间可能很长。

### 设置用户的默认shell
```
sudo sh -c 'echo /usr/bin/fish >> /etc/shells'
chsh -s /usr/bin/fish
```
