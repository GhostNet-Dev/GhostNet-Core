package fileservice

import (
	"log"
	"sync"

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
	FileObjs *sync.Map//map[string]*FileObject
}

func NewFileObjManager() *FileObjManager {
	return &FileObjManager{
		FileObjs: &sync.Map{},
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
	fileManager.FileObjs.Store(filename, fileObj)
	return fileObj
}

func (fileManager *FileObjManager) AllocBuffer(filename string, fileLength uint64) *FileObject {
	if fileObj, exist := fileManager.GetFileObject(filename); exist {
		fileObj.FileLength = fileLength
		fileObj.Buffer = make([]byte, fileLength)
		return fileObj
	}
	return nil
}

func (fileManager *FileObjManager) GetFileObject(filename string) (*FileObject, bool) {
	obj, ok := fileManager.FileObjs.Load(filename)
	return obj.(*FileObject), ok
}

func (fileManager *FileObjManager) DeleteObject(filename string) {
	fileManager.FileObjs.Delete(filename)
}

func (fileObj *FileObject) UpdateFileImage(offset uint64, bufSize uint64) bool {
	bitOffset := uint32(offset / bufSize)
	mutex := &sync.Mutex{}

	if fileObj.CompleteDone {
		return true
	}

	mutex.Lock()
	if fileObj.DownloadBitmap.Contains(bitOffset) {
		mutex.Unlock()
		return false
	}
	fileObj.DownloadBitmap.Set(bitOffset)
	mutex.Unlock()

	writtenSize := fileObj.FileLength - offset
	if writtenSize > bufSize {
		writtenSize = bufSize
	}

	fileObj.FileWrittenSize += writtenSize
	if fileObj.FileWrittenSize == fileObj.FileLength {
		fileObj.CompleteDone = true
		if fileObj.Callback != nil {
			fileObj.Callback(fileObj, fileObj.Context)
		}
	} else if fileObj.FileWrittenSize > fileObj.FileLength {
		log.Print("overflow WrittenSize=", fileObj.FileWrittenSize, ", FileLength=", fileObj.FileLength)
	}

	return true
}

func (fileObj *FileObject) CheckComplete() bool {
	return fileObj.CompleteDone
}

//func (fileManager *FileObjManager)
