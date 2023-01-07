package fileserver

type DoneHandler func(*FileObject, interface{})

type FileObject struct {
	Filename        string
	FileWrittenSize uint64
	FileLength      uint64
	Buffer          []byte
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
	callback DoneHandler, context interface{}) *FileObject {
	fileObj := &FileObject{
		Filename: filename,
		Buffer:   buffer,
		Callback: callback,
		Context:  context,
	}
	fileManager.FileObjs[filename] = fileObj
	return fileObj
}

func (fileManager *FileObjManager) GetFileObject(filename string) (*FileObject, bool) {
	obj, ok := fileManager.FileObjs[filename]
	return obj, ok
}
