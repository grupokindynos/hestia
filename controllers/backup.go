package controllers

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/base64"
	"encoding/json"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/aes"
	"github.com/grupokindynos/common/hestia"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/hestia/models"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"log"
	"os"
	"time"
)

type BackupController struct {
	Model *models.BackupModel
	Bucket string
}

func (bc *BackupController) CreateBackup(c *gin.Context) {
/*	_, _, err := mvt.VerifyRequest(c)
	if err != nil {
		responses.GlobalResponseNoAuth(c)
		return
	}*/

	_ = godotenv.Load()
	fbCredStr := os.Getenv("FIREBASE_CRED")
	fbCred, err := base64.StdEncoding.DecodeString(fbCredStr)
	if err != nil {
		log.Fatal("unable to decode firebase credential string:")
	}
	opt := option.WithCredentialsJSON(fbCred)
	fbApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal("unable to initialize firebase app")
	}

	storageClient, err := fbApp.Storage(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	storageBuck, err := storageClient.Bucket(bc.Bucket)
	if err != nil {
		log.Fatal(err)
	}

	query := &storage.Query{Prefix: ""}

	var lastBackup string
	var lastTime time.Time
	var size int64
	first := true
	it := storageBuck.Objects(context.Background(), query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if first {
			lastBackup = attrs.Name
			lastTime = attrs.Created
			size = attrs.Size
			first = false
		} else if attrs.Created.Unix() > lastTime.Unix() {
			lastBackup = attrs.Name
			lastTime = attrs.Created
			size = attrs.Size
		}
	}

	r, err := storageBuck.Object(lastBackup).NewReader(context.Background())
	if err != nil {
		log.Println(err)
	}
	defer r.Close()

	buf := make([]byte, size)

	if _, err := io.ReadFull(r, buf); err != nil {
		log.Println(err)
	}

	obj, err := aes.Decrypt([]byte(os.Getenv("HESTIA_BK_ENCRYPTION_KEY")), string(buf))
	if err != nil {
		log.Println(err)
	}

	objStr, err := base64.StdEncoding.DecodeString(obj)
	if err != nil {
		log.Println(err)
	}

	backupDb := hestia.HestiaDB{}
	err = json.Unmarshal(objStr, &backupDb)

	err = bc.Model.CreateBackup(backupDb)
	if err != nil {
		log.Println(err)
		responses.GlobalResponseError(nil, err, c)
		return
	}
}
