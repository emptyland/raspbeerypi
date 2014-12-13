package model

/*
#include <sys/vfs.h>
#include <stdio.h>

int bridge_statfs(const char *mp, int *entry) {
	struct statfs buf;

	int rv = statfs(mp, &buf);
	if (rv < 0)
		return rv;
	
	entry[0] = buf.f_bsize * buf.f_blocks;
	entry[1] = buf.f_bsize * (buf.f_blocks - buf.f_bfree);
	return 0;
}
*/
import "C"

import (
	"unsafe"
	"fmt"
	"bufio"
	"log"
	"os"
	"regexp"

	"api"
)

type DiskUsageModel struct {
}

var _ = (api.Model)(&DiskUsageModel{})

func (self *DiskUsageModel) Access(appKey string, token string) bool {
	return true
}

type mtabEntry struct {
	MountPoint string
	FSType     string
}

func (self *DiskUsageModel) GetApiDisk(res *diskUsageResponse) error {
	entries, err := readMtab()
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("mount points:", entries, err)

	res.Entries = make([]diskUsageVO, len(entries))
	var values [3]int

	for i, entry := range entries {
		code := C.bridge_statfs(C.CString(entry.MountPoint), (*C.int)(unsafe.Pointer(&values[0])))
		if code < 0 {
			log.Fatal("statfs error")
			return fmt.Errorf("statfs() fail")
		}

		res.Entries[i] = diskUsageVO {
			MountPoint: entry.MountPoint,
			FSType: entry.FSType,
			Total: int64(values[0]),
			Used:  int64(values[1]),
		}
	}
	return nil
}

var kSpaceRe = regexp.MustCompile("\\s+")
var kIgnoreFs = map[string]bool{
	"rootfs":     true,
	"devtmpfs":   true,
	"proc":       true,
	"sysfs":      true,
	"securityfs": true,
	"devpt":      true,
	"devpts":     true,
	"cgroup":     true,
	"pstore":     true,
	"systemd-1":  true,
	"mqueue":     true,
	"debugfs":    true,
	"binfmt_misc":true,
}

func readMtab() ([]mtabEntry, error) {

	mtabFile, err := os.Open(kMtabPath)
	if err != nil {
		return nil, err
	}
	defer mtabFile.Close()

	entries := make([]mtabEntry, 0)

	reader := bufio.NewReader(mtabFile)
	for {
		var raw []byte
		if raw, _, err = reader.ReadLine(); err != nil {
			break
		}

		// [0] device
		// [1] mount point
		// [2] fs type
		segments := kSpaceRe.Split(string(raw), 6)
		if kIgnoreFs[segments[0]] {
			continue
		}

		log.Println("mtab line:", segments)
		entries = append(entries, mtabEntry{
			MountPoint: segments[1],
			FSType:     segments[2],
		})
	}

	return entries, nil
}
