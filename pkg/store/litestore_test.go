package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

var (
	TestTables  = []string{DefaultNodeTable, DefaultWalletTable, "tdd"}
	db          = NewLiteStore("./", "litestore.db", TestTables)
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

type Key struct {
	time int64
	key  string
}

func TestSort(t *testing.T) {
	err := db.OpenStore()
	defer db.Close()
	assert.Equal(t, nil, err, "db open error")
	s := rand.NewSource(time.Now().Unix())
	randKey := rand.New(s)
	insertToDb(int64(randKey.Intn(100)), testAddress.GetPubAddress(), fmt.Sprint("inter", randKey.Intn(100)))

	for i := 0; i < 10; i++ {
		keyorder := int64(randKey.Intn(100))
		insertToDb(keyorder, testAddress.GetPubAddress(), fmt.Sprint("hellow-", keyorder))
	}

	fmt.Println("test result")
	keys, values, _ := db.LoadEntryDesc("tdd")
	for idx, key := range keys {
		fmt.Printf("key = %x, value = %s\n", key, values[idx])
	}
}

func insertToDb(sortKey int64, key string, value string) {
	entry := &Key{time: sortKey, key: key}
	byteKey, _ := marshal(entry)
	db.SaveEntry("tdd", byteKey, []byte(fmt.Sprint(value, entry.time)))
}

func marshal(data *Key) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, data.time); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, []byte(data.key)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func unmarshal(data interface{}, bs []byte) error {
	buf := bytes.NewBuffer(bs)
	if err := binary.Read(buf, binary.BigEndian, data); err != nil {
		return err
	}
	return nil
}
