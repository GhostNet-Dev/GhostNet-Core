package cloudservice

import (
	"fmt"
	"net"

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

	fileObj := <-cloud.streamId[filename]
	return fileObj
}

func (cloud *CloudService) DownloadDone(fileObj *fileservice.FileObject, context interface{}) {
	cloud.streamId[fileObj.Filename] <- fileObj
}
