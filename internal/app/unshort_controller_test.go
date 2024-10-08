package app

import (
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/shs"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_mainController_Get(t *testing.T) {
	type want struct {
		code        int
		response    string
		redirect    string
		contentType string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "positive test #1",
			request: "/qj",
			want: want{
				code:        307,
				redirect:    `https://practicum.yandex.ru/`,
				response:    ``,
				contentType: "",
			},
		},
		{
			name:    "negative test #2",
			request: "/",
			want: want{
				code:        400,
				redirect:    ``,
				response:    `Redirect URL not found`,
				contentType: "text/plain",
			},
		},
		{
			name:    "negative test #3",
			request: `/qweuoqiweuoiq_/*weu/kdalsdk;las?qweoipoq=73817`,
			want: want{
				code:        400,
				redirect:    ``,
				response:    `Redirect URL not found`,
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		unShort := createUnShortController()
		short := createShortController()
		e := Create().Build().GetEcho()
		t.Run(tt.name, func(t *testing.T) {
			//Пдготовка
			postReader := strings.NewReader(tt.want.redirect)
			postReq := httptest.NewRequest(http.MethodGet, "/", postReader)

			rec := httptest.NewRecorder()
			postContext := createContext(e, postReq, rec)
			postHandler := short.Post()
			postHandler(postContext)

			//Тест
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			c := createContext(e, request, w)

			err := unShort.Get()(c)

			if err != nil {
				assert.Error(t, err, tt.want.response)
			} else {
				res := w.Result()
				// проверяем код ответа
				assert.Equal(t, tt.want.code, res.StatusCode)
				// получаем и проверяем тело запроса
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				assert.Contains(t, string(resBody), tt.want.response)
				assert.Contains(t, res.Header.Get("Content-Type"), tt.want.contentType)
				assert.Equal(t, tt.want.redirect, res.Header.Get("Location"))
			}
		})
	}
}

func Test_mainController_Post(t *testing.T) {
	con := createUnShortController()
	assert.Nil(t, con.Post())
}

func createContext(e *echo.Echo, request *http.Request, w *httptest.ResponseRecorder) *RoteContext {
	ech := e.NewContext(request, w)
	return &RoteContext{
		ech,
		&storage.User{
			ID: 1,
		},
	}
}

func createUnShortController() *unShortController {
	conf := &config.AppConfig{RedirectAddress: url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	},
		HostName: "localhost:8080",
	}
	return &unShortController{
		path:    "/*",
		service: shs.NewShortenerService(storage.NewInMemoryRecordStorage(), storage.NewInMemoryUsersStorage(), conf),
		config:  conf,
	}
}
