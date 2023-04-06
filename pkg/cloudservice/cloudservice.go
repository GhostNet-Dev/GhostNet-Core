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

func (cloud *CloudService) DownloadASync(filename string, ipAddr *net.UDPAddr) <-chan *fileservice.FileObject {
	cloud.glog.DebugOutput(cloud, fmt.Sprint("cloud download = ", filename), glogger.Default)
	if _, exist := cloud.streamId[filename]; exist {
		return nil
	} else {
		cloud.streamId[filename] = nil
	}

	if ipAddr == nil {
		ipAddr = cloud.tTree.GetTreeClusterPick(filename)
	}
	cloud.fileService.SendGetFile(filename, ipAddr, cloud.DownloadDone, nil)

	go func() {
		defer cloud.ReleaseChannel(filename)
		select {
		case ret := <-cloud.streamId[filename]:
			if ret != nil {
				return
			}
		}
	}()
	return cloud.streamId[filename]
}

func (cloud *CloudService) DownloadDone(fileObj *fileservice.FileObject, context interface{}) {
	cloud.streamId[fileObj.Filename] <- fileObj
}
