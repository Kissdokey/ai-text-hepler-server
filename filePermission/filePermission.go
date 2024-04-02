package filePermission

import (
	"ai-text-helper-server/mysql"
)
const (
	Private = -1
	Visible = 0
	Commentable = 1
	Editable = 2
)

func CanISee(username string,fileId string) (bool,error) {
	fileInfo,err := mysql.GetFileInfo(fileId) 
	if err != nil {
		return false,err
	}
	if username == fileInfo.Author {
		return true,nil
	}
	if fileInfo.Permissions > Private {
		return true,nil
	}
	return false,nil

}

func CanIComment(username string,fileId string ) (bool,error) {
	fileInfo,err := mysql.GetFileInfo(fileId) 
	if err != nil {
		return false,err
	}
	if username == fileInfo.Author {
		return true,nil
	}
	if fileInfo.Permissions > Visible {
		return true,nil
	}
	return false,nil
}

func CanIEdit(username string,fileId string ) (bool,error) {
	fileInfo,err := mysql.GetFileInfo(fileId) 
	if err != nil {
		return false,err
	}
	if username == fileInfo.Author {
		return true,nil
	}
	if fileInfo.Permissions > Commentable {
		return true,nil
	}
	return false,nil
}
func GetPermission(username string,fileId string) int {
	fileInfo,err := mysql.GetFileInfo(fileId) 
	if err != nil {
		return Private
	}
	if username == fileInfo.Author {
		return Editable
	}
	//二期加白名单
	return fileInfo.Permissions
}