package variables

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// request body: { key: string, value: string }
func NewVariable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	user, err := api.NewContext(r).User()
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 64)
	if err != nil {
		api.NewResponse(w).Status(http.StatusBadRequest).Error(err.Error())
		return
	}

	var variable database.Variable
	if err := json.NewDecoder(r.Body).Decode(&variable); err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error(err.Error())
		return
	}
	variable.RepositoryID = repoID

	// SERIALIZE VARIABLE
	// gen id

	var wg sync.WaitGroup
	var serializeErr error
	go serializeVariable(&wg, &variable, serializeErr)

	// CHECK FOR WRITE ACCESS
	wa, err := db.hasWriteAccess(user.ID, repoID)
	if !wa {
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error("you do not have write access, please contact the repository owner.")
		return
	}

	// WRITE TO DB
	wg.Wait()

	if serializeErr != nil {
		api.NewResponse(w).ServerError(serializeErr.Error())
		return
	}

	// time
	variable.CreatedAt = database.TimeNow()
	variable.UpdatedAt = database.TimeNow()
	err = newVariable(&variable)
	if err == nil {
		api.NewResponse(w).Status(http.StatusCreated).JSON(&variable)
		return
	}

	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	log.Println(pgErr)
	if pgErr.Code == pgerrcode.UniqueViolation {
		api.NewResponse(w).
			Status(http.StatusConflict).
			Error(pgErr.Error())
		return
	}

	api.NewResponse(w).ServerError(pgErr.Error())
}

// writes new variable to db and update the `variable_count` of column of the same repo
func newVariable(v *database.Variable) error {
	return db.Transaction(func(tx *gorm.DB) error {
		defer func(tx *gorm.DB) {
			err, ok := recover().(string)
			if ok {
				tx.Rollback()
				log.Println(err)
			}
		}(tx)

		var err error = nil
		err = tx.Create(v).Error
		if err != nil {
			return err
		}

		var repo database.Repository
		err = tx.Where("id = ?", v.RepositoryID).First(&repo).Error
		if err != nil {
			return err
		}

		repo.VariableCount++
		return tx.Table(database.TableRepos).
			Where("id = ?", repo.ID).
			Updates(&repo).
			Error
	})
}
