package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/zakisanbaiman/go-handson01/testutil"
)

func TestLogin_ServeHTTP(t *testing.T) {
	type moq struct {
		token string
		err   error
	}
	type want struct {
		status  int
		repFile string
	}
	tests := map[string]struct {
		repFile string
		moq     moq
		want    want
	}{
		"ok": {
			repFile: "testdata/login/ok_req.json.golden",
			moq: moq{
				token: "from_moq",
			},
			want: want{
				status:  http.StatusOK,
				repFile: "testdata/login/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			repFile: "testdata/login/bad_req.json.golden",
			moq:     moq{},
			want: want{
				status:  http.StatusBadRequest,
				repFile: "testdata/login/bad_rsp.json.golden",
			},
		},
		"internalServerError": {
			repFile: "testdata/login/internal_server_error.json.golden",
			moq: moq{
				err: errors.New("error from mock"),
			},
			want: want{
				status:  http.StatusInternalServerError,
				repFile: "testdata/login/internal_server_error_rsp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		n := n
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/login",
				bytes.NewReader(testutil.LoadFile(t, tt.repFile)),
			)

			// mock
			moq := &LoginServiceMock{}
			moq.LoginFunc = func(ctx context.Context, name, pw string) (string, error) {
				return tt.moq.token, tt.moq.err
			}

			sut := Login{Service: moq, Validator: validator.New()}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t,
				resp,
				tt.want.status,
				testutil.LoadFile(t, tt.want.repFile),
			)
		})
	}
}
