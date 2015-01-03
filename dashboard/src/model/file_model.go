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
	fileModel := &FileManagementModel{
		fileEnv: fileEnv,
	}

	return api.NewService(fileModel)
}

// api/file/list
// api/file/delete
// api/file/new
// api/file/download/path/to/file
// api/file/upload/path/to/file
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
			vo.Type = "Directory"
		} else {
			vo.Type = "File"
		}

		rep.Entries = append(rep.Entries, vo)
	}

	api.EncodeJson(w, &rep)
	return nil
}
