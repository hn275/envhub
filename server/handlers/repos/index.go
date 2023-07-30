package repos

import (
	"encoding/json"
	"net/http"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/gh"
	"github.com/hn275/envhub/server/jsonwebtoken"
	"gorm.io/gorm"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.NewResponse(w).
			Status(http.StatusMethodNotAllowed).
			Error("http method not allowed")
		return
	}

	user, err := jsonwebtoken.GetUser(r)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error(err.Error())
		return
	}

	// GET REQUEST QUERY PARAM
	page := r.URL.Query().Get("page")
	sort := r.URL.Query().Get("sort")
	show := r.URL.Query().Get("show")
	if page == "" || sort == "" || show == "" {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("missing required queries")
		return
	}
	params := map[string]string{
		"page":     page,
		"sort":     sort,
		"per_page": show,
	}

	// NOTE: Since we are only interested in the repo that got sent back
	// by Github, this ops won't be a go routine.
	// get repos from github, then query db for the id of the same set of repos.

	// GET REPOS FROM GITHUB
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repositories-for-the-authenticated-user
	var repos []Repository
	ghCtx := gh.New(user.Token).Params(params)

	res, err := ghCtx.Get("/user/repos")
	defer res.Body.Close()
	if err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	if err := json.NewDecoder(res.Body).Decode(&repos); err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	ids := make([]uint64, len(repos))
	for i, repo := range repos {
		ids[i] = repo.ID
	}

	// GET REPO ID's FROM DB
	dbRepos, err := db.findRepo(user.ID, ids[:])
	switch err {
	case nil:
		isLinked := make([]bool, maxIDVal(repos)+1)
		for _, v := range dbRepos {
			isLinked[v] = true
		}

		for i := range repos {
			repos[i].Linked = isLinked[repos[i].ID]
		}
		break

	case gorm.ErrRecordNotFound:
		break

	default:
		api.NewResponse(w).ServerError(err)
		return
	}

	api.NewResponse(w).
		Header("Cache-Control", "max-age=30").
		Status(http.StatusOK).
		JSON(&repos)
}

func maxIDVal(ids []Repository) uint64 {
	var max uint64 = 0
	for _, v := range ids {
		if v.ID > max {
			max = v.ID
		}
	}

	return max
}
