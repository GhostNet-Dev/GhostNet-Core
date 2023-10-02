package cloudservice

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gnetwork"
)

type CloudService struct {
	fileService *fileservice.FileService
	tTree       *gnetwork.TrieTreeMap
	glog        *glogger.GLogger
	streamId    map[string](chan *fileservice.FileObject)
}

func NewCloudService(fileService *fileservice.FileService,
	tTree *gnetwork.TrieTreeMap, glog *glogger.GLogger) *CloudService {
	return &CloudService{
		fileService: fileService,
		tTree:       tTree,
		glog:        glog,
		streamId:    make(map[string]chan *fileservice.FileObject),
	}
}

func (cloud *CloudService) ReleaseChannel(filename string) {
	if _, exist := cloud.streamId[filename]; !exist {
		log.Print("already released stream = ", filename)
		return
	}
	close(cloud.streamId[filename])
	delete(cloud.streamId, filename)
}

func (cloud *CloudService) ReadFromCloudSync(filename string, ipAddr *net.UDPAddr) *fileservice.FileObject {
	defer cloud.ReleaseChannel(filename)
	cloud.glog.DebugOutput(cloud, fmt.Sprint("cloud download = ", filename), glogger.Default)
	if _, exist := cloud.streamId[filename]; !exist {
		cloud.streamId[filename] = make(chan *fileservice.FileObject, 1)
	}

	if ipAddr == nil {
		ipAddr = cloud.tTree.GetTreeClusterPick(filename)
	}

	cloud.fileService.SendGetFile(filename, ipAddr, cloud.DownloadDone, nil)

	select {
	case fileObj := <-cloud.streamId[filename]:
		return fileObj
	case <-time.After(time.Second * time.Duration(16)):
		cloud.glog.DebugOutput(cloud, fmt.Sprint("timeout download = ", filename), glogger.Default)
		return nil
	}
}

func (cloud *CloudService) DownloadDone(fileObj *fileservice.FileObject, context interface{}) {
	cloud.streamId[fileObj.Filename] <- fileObj
}
