package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/avast/retry-go"
	"github.com/labstack/echo/v4"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/storage"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/utils"
	"go.uber.org/zap"
)

func PostMetricsWithJSON(c echo.Context) error {
	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	var m = model.Metrics{}
	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	err := retry.Do(func() error {
		var err error
		m, err = storage.MStorage.SaveMetric(c.Request().Context(), m)
		if err != nil {
			if !utils.IsRetriableError(err) {
				return retry.Unrecoverable(err)
			}
			return err
		}
		return nil
	},
		retry.Attempts(4),
		retry.DelayType(utils.RetryDelay))
	if err != nil {
		logger.Log.Error("err", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, m)
}

func PostMetrics(c echo.Context) error {
	params := []string{
		c.Param("type"),
		c.Param("name"),
		c.Param("value"),
	}

	for _, param := range params {
		if param == "" {
			return c.String(http.StatusBadRequest, "invalid params")
		}
	}
	var m = model.Metrics{
		MType: params[0],
		ID:    params[1],
	}
	err := m.SetValue(params[2])
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = retry.Do(func() error {
		_, err = storage.MStorage.SaveMetric(c.Request().Context(), m)
		if err != nil {
			if !utils.IsRetriableError(err) {
				return retry.Unrecoverable(err)
			}
			return err
		}
		return nil
	},
		retry.Attempts(4),
		retry.DelayType(utils.RetryDelay))
	if err != nil {
		logger.Log.Error("err", zap.Error(err))
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "200 - OK")
}

func GetMetricIDWithJSON(c echo.Context) error {
	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}
	var m = model.Metrics{}
	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}
	err := retry.Do(func() error {
		var err error
		m, err = storage.MStorage.GetByID(c.Request().Context(), m.ID)
		if err != nil {
			if !utils.IsRetriableError(err) {
				return retry.Unrecoverable(err)
			}
			return err
		}
		return nil
	},
		retry.Attempts(4),
		retry.DelayType(utils.RetryDelay))

	if err != nil {
		logger.Log.Error("err", zap.Error(err))
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, m)
}

func GetMetricID(c echo.Context) error {
	params := []string{
		c.Param("type"),
		c.Param("name"),
	}

	for _, param := range params {
		if param == "" {
			return c.String(http.StatusNotFound, "invalid params")
		}
	}
	var m = model.Metrics{}
	err := retry.Do(func() error {
		var err error
		m, err = storage.MStorage.GetByID(c.Request().Context(), params[1])
		if err != nil {
			if !utils.IsRetriableError(err) {
				return retry.Unrecoverable(err)
			}
			return err
		}
		return nil
	},
		retry.Attempts(4),
		retry.DelayType(utils.RetryDelay))

	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.String(http.StatusOK, m.GetValueAsString())
}

func GetAllMetrics(c echo.Context) error {
	var metrics = map[string]model.Metrics{}

	err := retry.Do(func() error {
		var err error
		metrics, err = storage.MStorage.GetAllMetrics(c.Request().Context())
		if err != nil {
			if !utils.IsRetriableError(err) {
				return retry.Unrecoverable(err)
			}
			return err
		}
		return nil
	},
		retry.Attempts(4),
		retry.DelayType(utils.RetryDelay))

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	counters := make(map[string]model.Metrics)
	gauges := make(map[string]model.Metrics)

	for id, metric := range metrics {
		if metric.MType == "counter" {
			counters[id] = metric
		}
		if metric.MType == "gauge" {
			gauges[id] = metric
		}
	}

	tmp := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Metrics</title>
		</head>
		<body>
			<h1>Metrics</h1>

			<h2>Counters</h2>
			<table border="1">
				<tr>
					<th>Name</th>
					<th>Value</th>
				</tr>
				{{ range $name, $metric := .Counters }}
				<tr>
					<td>{{ $name }}</td>
					<td>{{ GetValueAsString $metric }}</td>
				</tr>
				{{ end }}
			</table>

			<h2>Gauges</h2>
			<table border="1">
				<tr>
					<th>Name</th>
					<th>Value</th>
				</tr>
				{{ range $name, $metric := .Gauges }}
				<tr>
					<td>{{ $name }}</td>
					<td>{{ GetValueAsString $metric }}</td>
				</tr>
				{{ end }}
			</table>

		</body>
		</html>
		`

	t, err := template.New("metrics").
		Funcs(template.FuncMap{"GetValueAsString": func(m model.Metrics) string {
			return m.GetValueAsString()
		},
		}).Parse(tmp)
	if err != nil {
		return c.String(http.StatusInternalServerError, "the template could not be executed")
	}

	data := struct {
		Counters map[string]model.Metrics
		Gauges   map[string]model.Metrics
	}{
		Counters: counters,
		Gauges:   gauges,
	}

	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")

	if err = t.Execute(c.Response().Writer, data); err != nil {
		return c.String(http.StatusInternalServerError, "the template could not be executed")
	}
	return nil
}

func PingHandler(c echo.Context) error {
	err := storage.MStorage.Ping(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to connect to the database")
	}
	return c.String(http.StatusOK, "Pong")
}

func UpdateSomeMetrics(c echo.Context) error {
	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	var metrics []model.Metrics
	if err := c.Bind(&metrics); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := retry.Do(func() error {
		var err error
		metrics, err = storage.MStorage.SaveMetrics(c.Request().Context(), metrics)
		if err != nil {
			if !utils.IsRetriableError(err) {
				return retry.Unrecoverable(err)
			}
			return err
		}
		return nil
	})

	if err != nil {
		logger.Log.Error("err", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, metrics)
}
