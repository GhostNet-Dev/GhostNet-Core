package cloudservice

import (
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gnetwork"
)

type CloudService struct {
	fileService *fileservice.FileService
	tTree       *gnetwork.TrieTreeMap
	streamId    map[string](chan *fileservice.FileObject)
}

func NewCloudService(fileService *fileservice.FileService,
	tTree *gnetwork.TrieTreeMap) *CloudService {
	return &CloudService{
		fileService: fileService,
		tTree:       tTree,
		streamId:    make(map[string]chan *fileservice.FileObject),
	}
}

func (cloud *CloudService) ReleaseChannel(filename string) {
	close(cloud.streamId[filename])
	delete(cloud.streamId, filename)
}

func (cloud *CloudService) DownloadASync(filename string, ipAddr *net.UDPAddr) <-chan *fileservice.FileObject {
	if _, exist := cloud.streamId[filename]; exist == true {
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
