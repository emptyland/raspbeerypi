
#include <stdio.h>
#include <string.h> 
#include <errno.h>
#include <string.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <sys/ioctl.h>
#include <net/if.h>

#include <wiringPi.h>
#include <lcd.h>

const int RS = 3;	// 
const int EN = 14;	// 
const int D0 = 4;	// 
const int D1 = 12;	// 
const int D2 = 13;	// 
const int D3 = 6;	// 

const char kNoIPv4[] = "---.---.---.---";

static int get_local_ipv4(const char *interface, char *ip, size_t n) {
	int fd = socket(PF_INET, SOCK_DGRAM, IPPROTO_IP);
	if (fd < 0) {
		perror("socket()");
		strncpy(ip, kNoIPv4, n);
		return fd;
	}

	struct ifreq req;
	strncpy(req.ifr_name, interface, IFNAMSIZ);
	if (ioctl(fd, SIOCGIFADDR, &req) < 0) {
		perror("ioctl()");
		strncpy(ip, kNoIPv4, n);
		close(fd);
		return -1;
	}

	struct sockaddr_in addr_in;
	memcpy(&addr_in, &req.ifr_addr, sizeof(addr_in));
	strncpy(ip, inet_ntoa(addr_in.sin_addr), n);
	close(fd);
	return 0;
}

static int get_cpu_temp(int *temp) {

	FILE *fp = fopen("/sys/class/thermal/thermal_zone0/temp", "r");

	if (fp == NULL) {
		perror("fopen()");
		return -1;
	}

	if (fscanf(fp, "%d", temp) <= 0) {
		perror("fscanf()");
		fclose(fp);
	}

	fclose(fp);
	return 0;
}

int main(int argc, char *argv[]) {
	const char *interface = "eth0";
	if (argc > 1)
		interface = argv[1];

	if (wiringPiSetup() < 0) {
		fprintf (stderr, "Unable to open serial device: %s\n", strerror (errno)) ;
		return 1 ;
	}

	int lcd = lcdInit(2, 16, 4, RS, EN, D0, D1, D2, D3, D0, D1, D2, D3);
	if (lcd < 0) {
		fprintf (stderr, "lcd init fail: %s\n", strerror (errno)) ;
		return 1 ;
	}

	char ipv4[16];
	get_local_ipv4(interface, ipv4, sizeof(ipv4));
	printf("%s\n", ipv4);
	lcdPosition(lcd, 0, 0);
	lcdPrintf(lcd, ipv4);

	lcdPosition(lcd, 0, 1);
	int temp;
	if (get_cpu_temp(&temp) < 0)
		lcdPrintf(lcd, " CPU: ---- C");
	else
		lcdPrintf(lcd, " CPU: %0.2f C", temp / 1000.0);
	return 0;
}
