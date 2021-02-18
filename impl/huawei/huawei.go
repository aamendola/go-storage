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
func MakeObjectStorage(endpoint, accessKey, secretAccessKey, bucketname, proxy string) ObjectStorage {

	var obsClient *obs.ObsClient
	var err error

	if len(endpoint) == 0 || len(accessKey) == 0 || len(secretAccessKey) == 0 || len(bucketname) == 0 {
		panic("ObjectStorage config error")
	}

	if proxy != "" {
		obsClient, err = obs.New(accessKey, secretAccessKey, endpoint, obs.WithProxyUrl(proxy))
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

// Get ...
func (o ObjectStorage) Get(filenameToDownload, filenameToSave string) {
	input := &obs.DownloadFileInput{}
	input.Bucket = o.bucketName
	input.Key = filenameToDownload
	input.DownloadFile = filenameToSave // Full path to which objects are downloaded.
	input.EnableCheckpoint = true       // Enable the resumable download mode.
	input.PartSize = 9 * 1024 * 1024    // Set the part size to 9 MB.
	input.TaskNum = 5                   // Specify the maximum number of parts that can be concurrently downloaded.
	output, err := o.client.DownloadFile(input)
	if err == nil {
		log.Printf("RequestId:%s\n", output.RequestId)
	} else if obsError, ok := err.(obs.ObsError); ok {
		log.Printf("Code:%s\n", obsError.Code)
		log.Printf("Message:%s\n", obsError.Message)
	}
}
