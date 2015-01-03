package model

import (
	"api"
	"log"
	"net/http"
	"os"
	"io"
)

type FileManagementModel struct {
	fileEnv *FileEnvDef
}

var _ = (api.Model)(&FileManagementModel{})
var _ = log.Println

func NewFileManagementService(fileEnv *FileEnvDef) *api.Service {
	if fileEnv.Home == "" {
		log.Fatal("file_env.home can not be empty!")
	}

	fileModel := &FileManagementModel{
		fileEnv: fileEnv,
	}

	return api.NewService(fileModel)
}

// api/file/{list|delete|new|download|upload}/path/to/file
func (self *FileManagementModel) Access(appKey string, token string) bool {
	return true
}

func (self *FileManagementModel) PrefixGetApiFileDownload(w http.ResponseWriter, r *http.Request, path string) error {
	fileName := self.fileEnv.Home + path

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	return err
}

func (self *FileManagementModel) PrefixGetApiFileList(w http.ResponseWriter, r *http.Request, path string) error {
	rep := fileStatReponse{
		Entries: make([]fileStatVO, 0),
	}

	dir, err := os.Open(self.fileEnv.Home + path)
	if err != nil {
		return err
	}
	defer dir.Close()

	var fileInfos []os.FileInfo
	if fileInfos, err = dir.Readdir(0); err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		switch fileInfo.Name() {
		case ".", "..":
			continue
		}

		vo := fileStatVO{
			Name:    fileInfo.Name(),
			Size:    fileInfo.Size(),
			Mode:    int16(fileInfo.Mode()),
			ModTime: fileInfo.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:   fileInfo.IsDir(),
		}

		if fileInfo.IsDir() {
			vo.Type = "Directory";
			if vo.Contain, err = getDirContain(self.fileEnv.Home + path + "/" + vo.Name); err != nil {
				return err
			}
		} else {
			vo.Type = "File";
			vo.Contain = -1;
		}

		rep.Entries = append(rep.Entries, vo)
	}

	api.EncodeJson(w, &rep)
	return nil
}

func getDirContain(dirName string) (int, error) {
	dir, err := os.Open(dirName)
	if err != nil {
		return 0, err
	}
	defer dir.Close()

	var fileInfos []os.FileInfo
	if fileInfos, err = dir.Readdir(0); err != nil {
		return 0, err
	}

	return len(fileInfos), nil
}
