package reader

import (
	"encoding/csv"
	"github.com/lf-edge/ekuiper/internal/io/file/common"
	"io"
	"strconv"
	"strings"

	"github.com/lf-edge/ekuiper/pkg/api"
)

type CsvReader struct {
	csvR   *csv.Reader
	config *common.FileSourceConfig

	ctx api.StreamContext
}

func (r *CsvReader) Read() ([]map[string]interface{}, error) {
	record, err := r.csvR.Read()
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		r.ctx.GetLogger().Warnf("Read file %s encounter error: %v", "fs.file", err)
		return nil, err
	}
	r.ctx.GetLogger().Debugf("Read" + strings.Join(record, ","))
	cols := r.config.Columns
	var m map[string]interface{}
	if cols == nil {
		m = make(map[string]interface{}, len(record))
		for i, v := range record {
			m["cols"+strconv.Itoa(i)] = v
		}
	} else {
		m = make(map[string]interface{}, len(cols))
		for i, v := range cols {
			m[v] = record[i]
		}
	}

	return []map[string]interface{}{m}, nil
}

func (r *CsvReader) Close() error {
	return nil
}

func CreateCsvReader(fileStream io.Reader, config *common.FileSourceConfig, ctx api.StreamContext) (FormatReader, error) {
	r := csv.NewReader(fileStream)
	r.Comma = rune(config.Delimiter[0])
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = -1
	cols := config.Columns
	if config.HasHeader {
		var err error
		ctx.GetLogger().Debug("Has header")
		cols, err = r.Read()
		if err == io.EOF {
			return nil, err
		}
		if err != nil {
			ctx.GetLogger().Warnf("Read file %s encounter error: %v", "fs.file", err)
			return nil, err
		}
		ctx.GetLogger().Debugf("Got header %v", cols)
	}

	reader := &CsvReader{}
	reader.csvR = r
	reader.config = config
	reader.ctx = ctx

	return reader, nil
}
