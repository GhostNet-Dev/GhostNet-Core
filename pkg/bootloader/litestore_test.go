package bootloader

import (
	"bytes"
	"log"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

var (
	TestTables  = []string{"nodes", "wallet"}
	db          = NewLiteStore("./", TestTables)
	testAddress = gcrypto.GenerateKeyPair()
	ghostIp     = &ptypes.GhostIp{
		Ip:   "127.0.0.1",
		Port: "8888",
	}
)

func TestInAndOut(t *testing.T) {
	err := db.OpenStore()
	defer db.Close()
	assert.Equal(t, nil, err, "db open error")
	testUser := &ptypes.GhostUser{
		PubKey:   testAddress.GetPubAddress(),
		Nickname: "test",
		Ip:       ghostIp,
	}
	rawData, err := proto.Marshal(testUser)
	if err != nil {
		log.Fatal(err)
	}
	db.SaveEntry(TestTables[0], []byte(testUser.PubKey), rawData)
	v, err := db.SelectEntry(TestTables[0], []byte(testUser.PubKey))
	assert.Equal(t, nil, err, "db select error")
	assert.Equal(t, 0, bytes.Compare(rawData, v), "load fail")

	db.LoadEntry(TestTables[0])
}
