package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/grupokindynos/common/hestia"
	"time"
)

type BackupModel struct {
	Firestore   *firestore.DocumentRef
	Document    string
}

func (bm *BackupModel) CreateBackup(hestiaDb hestia.HestiaDB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50 * time.Second)
	defer cancel()
	_, err := bm.Firestore.Set(ctx, hestiaDb)
	return err
}