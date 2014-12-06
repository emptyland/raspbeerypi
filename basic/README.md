

> 假设使用Pidora

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
