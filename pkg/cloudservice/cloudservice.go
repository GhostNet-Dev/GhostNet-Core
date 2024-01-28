package cloudservice

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gnetwork"
)

type CloudService struct {
	fileService *fileservice.FileService
	tTree       *gnetwork.TrieTreeMap
	glog        *glogger.GLogger
	streamId    *sync.Map //map[string](chan *fileservice.FileObject)
}

func NewCloudService(fileService *fileservice.FileService,
	tTree *gnetwork.TrieTreeMap, glog *glogger.GLogger) *CloudService {
	return &CloudService{
		fileService: fileService,
		tTree:       tTree,
		glog:        glog,
		streamId:    &sync.Map{}, //make(map[string]chan *fileservice.FileObject),
	}
}

func (cloud *CloudService) ReleaseChannel(filename string) {
	ch, exist := cloud.streamId.Load(filename)
	if !exist {
		log.Print("already released stream = ", filename)
		return
	}
	close(ch.(chan *fileservice.FileObject))
	cloud.streamId.Delete(filename)
}

func (cloud *CloudService) ReadFromCloudSync(filename string, ipAddr *net.UDPAddr) *fileservice.FileObject {
	defer cloud.ReleaseChannel(filename)
	cloud.glog.DebugOutput(cloud, fmt.Sprint("cloud download = ", filename), glogger.Default)

	ch, exist := cloud.streamId.Load(filename)
	if !exist {
		ch = make(chan *fileservice.FileObject, 1)
		cloud.streamId.Store(filename, ch)
	}

	if ipAddr == nil {
		ipAddr = cloud.tTree.GetTreeClusterPick(filename)
	}

	cloud.fileService.SendGetFile(filename, ipAddr, cloud.DownloadDone, nil)

	select {
	case fileObj := <-ch.(chan *fileservice.FileObject):
		return fileObj
	case <-time.After(time.Second * time.Duration(16)):
		cloud.glog.DebugOutput(cloud, fmt.Sprint("timeout download = ", filename), glogger.Default)
		return nil
	}
}

func (cloud *CloudService) DownloadDone(fileObj *fileservice.FileObject, context interface{}) {
	if ch, exist := cloud.streamId.Load(fileObj.Filename); exist {
		ch.(chan *fileservice.FileObject) <- fileObj
	} else {
		cloud.glog.DebugOutput(cloud, fmt.Sprint("download done miss = ", fileObj.Filename), glogger.Default)
	}
}
