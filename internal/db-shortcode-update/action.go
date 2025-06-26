package dbShortcodeUpdate

import (
	"errors"
	"fmt"
	fireStore "github.com/driscollco-core/firestore"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
)

func Action(shortCode entities.ShortCode, db fireStore.Client, values map[string]fireStore.UpdateValue) error {
	if err := db.Update(fmt.Sprintf("%s/%s", conf.Config.GCP.FireStore.Paths.ShortCodes, shortCode.Id), values); err != nil {
		if !errors.Is(err, fireStore.ErrorNoResults) {
			return err
		}
		if err = db.Update(fmt.Sprintf("%s/%s", conf.Config.GCP.FireStore.Paths.ShortCodesLegacy, shortCode.Id), values); err != nil {
			return err
		}
	}
	return nil
}
