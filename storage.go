package storage

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aamendola/go-storage/obs"
)

func init() {
	log.Printf("#####################################################################\n")
	log.Printf("#####################################################################\n")
}

// ObjectStorage is a configuration about Huawei OBS
type ObjectStorage struct {
	Client     *obs.ObsClient
	BucketName string
	Upload     bool
	Download   bool
}

// GetObjectStorageConfig ...
func NewObjectStorage(endpoint, ak, sk, bucketname string, uploadEnabled, downloadEnabled, proxyEnabled bool) (*ObjectStorage, error) {

	var obsClient *obs.ObsClient
	var err error

	if proxyEnabled {
		obsClient, err = obs.New(ak, sk, endpoint, obs.WithProxyUrl("http://proxy.mpba.gov.ar:3128"))
	} else {
		obsClient, err = obs.New(ak, sk, endpoint)
	}

	if err != nil {
		panic(err)
	}

	if !uploadEnabled || (len(endpoint) == 0 || len(ak) == 0 || len(sk) == 0 || len(bucketname) == 0) {
		return nil, fmt.Errorf("ObjectStorage config error")
	} else {
		fmt.Printf("ObjectStorage no habilitado")
	}

	objectStorage := ObjectStorage{obsClient, bucketname, uploadEnabled, downloadEnabled}
	return &objectStorage, nil
}

// Post ...
func (os *ObjectStorage) Post(src string, dst string) {

	data, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Printf("Error al leer archivo para generar md5: %s", err)
		panic(err)
	}

	base64Md5 := obs.Base64Md5(data)

	putFileInput := &obs.PutFileInput{}
	putFileInput.Bucket = os.BucketName
	putFileInput.Key = dst
	putFileInput.SourceFile = src
	putFileInput.ContentMD5 = base64Md5

	output, err := os.Client.PutFile(putFileInput)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("ETag:%s, StorageClass:%s\n", output.ETag, output.StorageClass)
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Printf("StatusCode:%d , Message:%s\n", output.StatusCode, obsError.Message)
		} else {
			fmt.Printf("Error:%s", err)
		}
	}

}

// Get ...
func (os *ObjectStorage) Get() {
}
