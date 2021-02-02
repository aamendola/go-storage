package huawei

import (
	"io/ioutil"
	"log"

	"github.com/aamendola/go-storage/impl/huawei/obs"
)

// ObjectStorage ...
type ObjectStorage struct {
	client     *obs.ObsClient
	bucketName string
}

// MakeObjectStorage ...
func MakeObjectStorage(endpoint, accessKey, secretAccessKey, bucketname string, proxy bool) ObjectStorage {

	var obsClient *obs.ObsClient
	var err error

	if len(endpoint) == 0 || len(accessKey) == 0 || len(secretAccessKey) == 0 || len(bucketname) == 0 {
		panic("ObjectStorage config error")
	}

	if proxy {
		obsClient, err = obs.New(accessKey, secretAccessKey, endpoint, obs.WithProxyUrl("http://proxy.mpba.gov.ar:3128"))
	} else {
		obsClient, err = obs.New(accessKey, secretAccessKey, endpoint)
	}

	if err != nil {
		panic(err)
	}

	return ObjectStorage{obsClient, bucketname}
}

// Post ...
func (os ObjectStorage) Post(filenameSource, filenameDestination string) {

	data, err := ioutil.ReadFile(filenameSource)
	if err != nil {
		log.Printf("Error al leer archivo para generar md5: %s", err)
		panic(err)
	}

	base64Md5 := obs.Base64Md5(data)

	putFileInput := &obs.PutFileInput{}
	putFileInput.Bucket = os.bucketName
	putFileInput.Key = filenameDestination
	putFileInput.SourceFile = filenameSource
	putFileInput.ContentMD5 = base64Md5

	output, err := os.client.PutFile(putFileInput)
	if err == nil {
		log.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		log.Printf("ETag:%s, StorageClass:%s\n", output.ETag, output.StorageClass)
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			log.Printf("StatusCode:%d , Message:%s\n", output.StatusCode, obsError.Message)
		} else {
			log.Printf("Error:%s", err)
		}
	}

}
