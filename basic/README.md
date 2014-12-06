

> 假设使用Pidora

## How to install fish shell

### install g++
```
sudo yum install gcc-c++
```

### install curses, dep by fish
```
sudo yum install ncurses-devel
```

### get the fish code
```
cd ~
wget 'http://fishshell.com/files/2.1.1/fish-2.1.1.tar.gz'
```

### compiler and install
```
tar vxzf fish-2.1.1.tar.gz
cd fish-2.1.1
./configure --prefix=/usr
make
sudo make install
```

### change user default shell
```
chsh -s /usr/bin/fish
```
