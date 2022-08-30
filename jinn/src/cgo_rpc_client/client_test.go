package client_test

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"jinnDataProcessing/config"
	rpcCli "jinnDataProcessing/jinn/src/cgo_rpc_client"
	"jinnDataProcessing/jinn/src/storage/define"
	"jinnDataProcessing/jinn/src/utils/log"
	"os"
	"testing"
)

func TestClient_DownloadOne(t *testing.T) {
	client := rpcCli.Client{User: "233686", Pwd: "aini1996.."}
	client.Connect("10.30.113.141", 10289)
	defer client.Disconnect()

	handle, err := client.FileOpen(1000, 2485986994308521858, "1025_2021-05-02120909__浙AF93051_新能源汽车.jpg")
	if err != nil {
		log.Fatal("%s", err)
	}
	h := sha1.New()
	// buf := make([]byte, 1024*256)
	offset := int64(0)
	for offset < 24000 {
		data, err := client.FileRead(handle, 1024*256)
		if err != nil {
			if e, ok := err.(define.Error); ok {
				if e.Code == define.ECEndOfFile {
					break
				}
			}
			log.Fatal("%s", err)
		}
		offset += int64(len(data))
		h.Write(data)
	}
	if err := client.FileClose(handle); err != nil {
		log.Fatal("%s", err)
	}
	s := hex.EncodeToString(h.Sum(nil))
	t.Log(s)
}

func TestClient_Upload(t *testing.T) {
	client := rpcCli.Client{User: "ucroot", Pwd: "111111"}
	client.Connect("10.30.113.141", 10388)
	defer client.Disconnect()
	file, err := os.Open(`C:\Users\278863\Pictures\20220307164245.png`)
	if err != nil {
		log.Fatal("%s", err)
	}
	info, err := file.Stat()
	if err != nil {
		log.Fatal("%s", err)
	}
	data := make([]byte, 512*1024)
	handle, err := client.UploadCreate(1000, 432345564227568733, info.Name(), info.Size())
	if err != nil {
		log.Fatal("%s", err)
	}
	defer func() {
		_ = file.Close()
		handle.Release()
	}()
	for {
		n, err := file.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal("%s", err)
			}
		}
		err = client.UploadWrite(handle, data[:n])
		if err != nil {
			log.Fatal("%s", err)
		}
	}
}

func TestClient_Dir(t *testing.T) {
	client := rpcCli.Client{User: config.JinnConfig.Username, Pwd: config.JinnConfig.Password}
	client.Connect(config.JinnConfig.Host, config.JinnConfig.Port)
	defer client.Disconnect()

	spaceId, err := client.GetSpaceId("test")

	pathInfo, err := client.GetPathInfo(spaceId, "/233686/pic")
	if err != nil {
		log.Fatal("%s", err)
	}

	h, err := client.DirOpen(spaceId, pathInfo.Id)
	if err != nil {
		log.Fatal("%s", err)
	}
	for {
		fis, err := client.DirRead(h, 1000)
		if err != nil {
			if e, ok := err.(define.Error); ok {
				if e.Code == define.ECEndOfFile {
					break
				}
			}
			log.Fatal("%s", err)
		}
		for _, i := range fis {
			log.Info("Name=%s, BaseID=%d", i.Name, i.BaseId)
		}
	}
	_ = client.DirClose(h)

}

func TestClient_GetTags(t *testing.T) {
	client := rpcCli.Client{User: "233686", Pwd: "aini1996.."}
	client.Connect("10.30.113.141", 10289)
	defer client.Disconnect()
	infos, err := client.GetTags(0, 1000, []int64{288230376151712783, 3458764513820541928})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", infos)
}

func TestClient_GetFileObjs(t *testing.T) {
	client := rpcCli.Client{User: "233686", Pwd: "aini1996.."}
	client.Connect("10.30.113.141", 10289)
	defer client.Disconnect()
	infos, err := client.GetFileObjs(666, []int64{1000, 2305843009213694952})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", infos)
}

func TestClient_GetFrameObjs(t *testing.T) {
	//file := database.FileInfo{BaseId: 3213213113333535, Class: "温度计阿达大大啪嗒啪嗒的爬坡打破打破阿婆打破"}
	//fmt.Println(unsafe.Sizeof(file))
	client := rpcCli.Client{User: "233686", Pwd: "aini1996.."}
	client.Connect("10.30.113.141", 10289)
	defer client.Disconnect()
	infos, err := client.GetFrameObjs(666, 4611686018427388904, [2]int32{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", infos)
}
