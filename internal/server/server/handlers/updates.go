package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"ya-metrics/pkg/mdata"
)

func (h *Handler) Updates(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var metrics []mdata.Metrics
	body, err := io.ReadAll(req.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &metrics)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, fmt.Sprintf("could not unmarshal data. err:%s", err), http.StatusBadRequest)
		return
	}
	if len(metrics) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, "empty metrics", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	for _, m := range metrics {
		var mType string
		/*		var delta sql.NullInt64
				var value sql.NullFloat64*/
		var delta int64
		var value float64
		switch strings.ToLower(m.MType) {
		case mdata.GAUGE:
			mType = mdata.GAUGE
			value = *m.Value
		case mdata.COUNTER:
			mType = mdata.COUNTER
			oldVal, err := h.countStorage.Get(m.ID)
			if err != nil {
				delta = *m.Delta
			} else {
				delta = oldVal.GetValue() + *m.Delta
			}
		default:
			h.log.Errorf("get undefened type:%s, ID:%s, Delta:%d, Value: %d", m.MType, m.ID, m.Delta, m.Value)
			continue
		}
		_, err = tx.ExecContext(req.Context(),
			"INSERT INTO metrics (id, mtype, delta, value) VALUES ($1,$2,$3,$4) "+
				"ON CONFLICT (id) DO UPDATE SET delta=EXCLUDED.delta, value=EXCLUDED.value",
			m.ID, mType, delta, value,
		)
		fmt.Println(m.ID, mType, delta, value)
		if err != nil {
			h.log.Errorf("error on transaction. err:%s", err)
			tx.Rollback()
			writer.WriteHeader(http.StatusInternalServerError)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		h.log.Errorf("tx commit err:%s", err)
	}
}
