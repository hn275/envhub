package variables

//
// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
//
// 	"github.com/go-chi/chi/v5"
// 	"github.com/hn275/envhub/server/database"
// 	"github.com/hn275/envhub/server/envhubtest"
// 	jwt "github.com/hn275/envhub/server/jsonwebtoken"
// 	"github.com/stretchr/testify/assert"
// )
//
// type mockGhCtxOK struct{}
// type mockGhCtxNotFound struct{}
// type mockGhCtxError struct{}
//
// var mockVar = database.Variable{
// 	Key:   "test_foo",
// 	Value: "test_bar",
// }
//
// func testInit(url, repoID string) (*httptest.ResponseRecorder, error) {
// 	buf, err := json.Marshal(mockVar)
// 	if err != nil {
// 		return nil, err
// 	}
// 	body := bytes.NewReader(buf)
//
// 	r := envhubtest.RequestWithParam(
// 		http.MethodPost,
// 		url,
// 		map[string]string{"repoID": repoID},
// 		body,
// 	)
// 	r.Header.Add("Content-Type", "application/json")
// 	r.Header.Add("Authorization", "Bearer "+"somejwttoken")
//
// 	m := chi.NewMux()
// 	m.Use(WriteAccessChecker)
// 	m.Handle("/repos/{repoID}/variables/new", http.HandlerFunc(NewVariable))
//
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, r)
// 	return w, nil
// }
//
// func cleanup() {
// 	d := database.New()
// 	d.Where("key = ?", "test_foo").Delete(&database.Variable{})
// }
//
// func TestNewVariable(t *testing.T) {
// 	defer cleanup()
// 	jwt.Mock()
//
// 	w, err := testInit("/repos/1/variables/new", "1")
// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
//
// 	var variable database.Variable
// 	d := database.New()
// 	err = d.Where("repository_id = ? AND key = ?", 1, "foo").First(&variable).Error
// 	assert.Nil(t, err)
// 	assert.NotEqual(t, mockVar.Value, variable.Value)
// }
//
// func TestNewVariableDuplicate(t *testing.T) {
// 	defer cleanup()
// 	jwt.Mock()
//
// 	w, err := testInit("/repos/1/variables/new", "1")
// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
//
// 	w, err = testInit("/repos/1/variables/new", "1")
// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusConflict, w.Result().StatusCode)
// }
//
// func TestWriteAccess(t *testing.T) {
// 	jwt.Mock()
//
// 	w, err := testInit("/repos/420/variables/new", "420")
// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
//
// 	// testing write-access ok
// 	defer cleanup()
// 	w, err = testInit("/repos/1/variables/new", "1")
// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
//
// 	var variable database.Variable
// 	d := database.New()
// 	err = d.Where("repository_id = ? AND key = ?", 1, "foo").First(&variable).Error
// 	assert.Nil(t, err)
// 	assert.NotEqual(t, mockVar.Value, variable.Value)
// }
//
// func TestInvalidRepoID(t *testing.T) {
// 	jwt.Mock()
//
// 	w, err := testInit("/repos/lkasjdf/variables/new", "ajsdf")
// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
// }
//
// func TestMethodNotAllowed(t *testing.T) {
// 	srv := httptest.NewServer(http.HandlerFunc(NewVariable))
// 	cx := http.Client{}
//
// 	methods := []string{
// 		http.MethodGet,
// 		http.MethodPut,
// 		http.MethodPatch,
// 		http.MethodHead,
// 		http.MethodDelete,
// 	}
// 	for _, method := range methods {
// 		req, err := http.NewRequest(method, srv.URL, nil)
// 		assert.Nil(t, err)
// 		res, err := cx.Do(req)
// 		assert.Nil(t, err)
// 		assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
// 	}
// }
//
// // MOCK
// func (m *mockGhCtxOK) Do(r *http.Request) (*http.Response, error) {
// 	res := &http.Response{
// 		Status:     "204 No Content",
// 		StatusCode: http.StatusNoContent,
// 		Request:    r,
// 	}
// 	return res, nil
// }
//
// func (m *mockGhCtxNotFound) Do(r *http.Request) (*http.Response, error) {
// 	res := &http.Response{
// 		Status:     "404 Not Found",
// 		StatusCode: http.StatusNotFound,
// 		Request:    r,
// 	}
// 	return res, nil
// }
//
// func (m *mockGhCtxError) Do(r *http.Request) (*http.Response, error) {
// 	res := &http.Response{
// 		Status:     "500 Internal Server Error",
// 		StatusCode: http.StatusInternalServerError,
// 		Request:    r,
// 	}
// 	return res, nil
// }
