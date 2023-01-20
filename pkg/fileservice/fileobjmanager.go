package fileservice

import (
	"log"

	"github.com/kelindar/bitmap"
)

type DoneHandler func(*FileObject, interface{})

type FileObject struct {
	Filename        string
	FileWrittenSize uint64
	FileLength      uint64
	Buffer          []byte
	DownloadBitmap  bitmap.Bitmap
	CompleteDone    bool
	Ratio           uint64
	Callback        DoneHandler
	Context         interface{}
}

type FileObjManager struct {
	FileObjs map[string]*FileObject
}

func NewFileObjManager() *FileObjManager {
	return &FileObjManager{
		FileObjs: make(map[string]*FileObject),
	}
}

func (fileManager *FileObjManager) CreateFileObj(filename string, buffer []byte,
	fileLength uint64, callback DoneHandler, context interface{}) *FileObject {
	completeDone := false
	fileWrittenSize := uint64(0)
	if buffer == nil {
		buffer = make([]byte, fileLength)
	} else {
		completeDone = true
		fileWrittenSize = fileLength
	}

	fileObj := &FileObject{
		Filename:        filename,
		FileLength:      fileLength,
		FileWrittenSize: fileWrittenSize,
		Buffer:          buffer,
		Callback:        callback,
		Context:         context,
		DownloadBitmap:  bitmap.Bitmap{},
		CompleteDone:    completeDone,
	}
	fileManager.FileObjs[filename] = fileObj
	return fileObj
}

func (fileManager *FileObjManager) GetFileObject(filename string) (*FileObject, bool) {
	obj, ok := fileManager.FileObjs[filename]
	return obj, ok
}

func (fileObj *FileObject) UpdateFileImage(offset uint64) bool {
	bitOffset := uint32(offset / BufferSize)
	if fileObj.DownloadBitmap.Contains(bitOffset) == true {
		return false
	}
	fileObj.DownloadBitmap.Set(bitOffset)
	writtenSize := fileObj.FileLength - offset
	if writtenSize > BufferSize {
		writtenSize = BufferSize
	}

	fileObj.FileWrittenSize += writtenSize
	if fileObj.FileWrittenSize == fileObj.FileLength {
		fileObj.CompleteDone = true
		if fileObj.Callback != nil {
			fileObj.Callback(fileObj, fileObj.Context)
		}
	} else if fileObj.FileWrittenSize > fileObj.FileLength {
		log.Fatal("overflow")
	}

	return true
}

func (fileObj *FileObject) CheckComplete() bool {
	return fileObj.CompleteDone
}

//func (fileManager *FileObjManager)
