package api

import (
	"github.com/IgorViskov/go_33_shortener/internal/config"
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
				response:    `<a href="https://practicum.yandex.ru/">Temporary Redirect</a>.` + "\n\n",
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name:    "negative test #2",
			request: "/",
			want: want{
				code:        400,
				redirect:    ``,
				response:    `Redirect URL not found`,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "negative test #3",
			request: `/qweuoqiweuoiq_/*weu/kdalsdk;las?qweoipoq=73817`,
			want: want{
				code:        400,
				redirect:    ``,
				response:    `Redirect URL not found`,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		c := NewMainController(&config.AppConfig{
			RedirectAddress: &url.URL{
				Scheme: "http",
				Host:   "localhost:8080",
			},
			BaseAddress: "localhost:8080",
		})
		t.Run(tt.name, func(t *testing.T) {
			//Пдготовка
			postReader := strings.NewReader(tt.want.redirect)
			postReq := httptest.NewRequest(http.MethodGet, "/", postReader)

			pw := httptest.NewRecorder()
			c.Post(pw, postReq)

			//Тест
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			c.Get(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Contains(t, string(resBody), tt.want.response)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.redirect, res.Header.Get("Location"))
		})
	}
}

func Test_mainController_Post(t *testing.T) {
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
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive test #2",
			body: `http://localhost:8080`,
			want: want{
				code:        201,
				response:    `http://localhost:8080/qj`,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #3",
			body: `/qweuoqiweuoiq_/*weu/kdalsdk;las?qweoipoq=73817`,
			want: want{
				code:        400,
				response:    `Invalid URL`,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #4",
			body: ``,
			want: want{
				code:        400,
				response:    `Invalid URL`,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		c := NewMainController(&config.AppConfig{
			RedirectAddress: &url.URL{
				Scheme: "http",
				Host:   "localhost:8080",
			},
			BaseAddress: "localhost:8080",
		})
		t.Run(tt.name, func(t *testing.T) {
			//Тест
			request := httptest.NewRequest(http.MethodGet, "localhost:8080", strings.NewReader(tt.body))
			// создаём новый Recorder
			w := httptest.NewRecorder()

			c.Post(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Contains(t, string(resBody), tt.want.response)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
