package api

import (
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/shs"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_shortController_Get(t *testing.T) {
	con := createShortController()
	assert.Nil(t, con.Get())
}

func Test_shortController_Post(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "positive test #1",
			body: `https://practicum.yandex.ru/`,
			want: want{
				code:        201,
				response:    `http://localhost:8080/qj`,
				contentType: "text/plain",
			},
		},
		{
			name: "positive test #2",
			body: `http://localhost:8080`,
			want: want{
				code:        201,
				response:    `http://localhost:8080/qj`,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #3",
			body: `/qweuoqiweuoiq_/*weu/kdalsdk;las?qweoipoq=73817`,
			want: want{
				code:        400,
				response:    `Invalid URL`,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #4",
			body: ``,
			want: want{
				code:        400,
				response:    `Invalid URL`,
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		con := createShortController()
		e := app.Create().Build().GetEcho()
		t.Run(tt.name, func(t *testing.T) {
			//Тест
			request := httptest.NewRequest(http.MethodGet, "localhost:8080", strings.NewReader(tt.body))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)

			var err = con.Post()(c)
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
			}
		})
	}
}

func createShortController() *shortController {
	return &shortController{
		path:    "/*",
		service: shs.NewShortenerService(storage.NewInMemoryStorage()),
		config: &config.AppConfig{RedirectAddress: url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
		},
			HostName: "localhost:8080",
		},
	}
}
