package merger_test

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"

	prom "github.com/prometheus/client_model/go"

	"github.com/vadv/prometheus-exporter-merger/merger"
)

const (
	data1 = `
# TYPE fluentbit_filter_add_records_total counter
fluentbit_filter_add_records_total{name="kubernetes.0"} 0 1589716338403
fluentbit_filter_add_records_total{name="lua.1"} 0 1589716338403
fluentbit_filter_add_records_total{name="lua.2"} 0 1589716338403
fluentbit_filter_add_records_total{name="lua.3"} 0 1589716338403
fluentbit_filter_add_records_total{name="kubernetes.0"} 0 1589716338417
`
	data2 = `
# TYPE fluentbit_filter_add_records_total counter
fluentbit_filter_add_records_total{name="lua.1"} 0 1589716338417
fluentbit_filter_add_records_total{name="rewrite_tag.2"} 0 1589716338417
`
	result = `
# TYPE fluentbit_filter_add_records_total counter
fluentbit_filter_add_records_total{name="kubernetes.0",url="value1"} 0 1589716338403
fluentbit_filter_add_records_total{name="kubernetes.0",url="value1"} 0 1589716338417
fluentbit_filter_add_records_total{name="lua.1",url="value1"} 0 1589716338403
fluentbit_filter_add_records_total{name="lua.1",url="value2"} 0 1589716338417
fluentbit_filter_add_records_total{name="lua.2",url="value1"} 0 1589716338403
fluentbit_filter_add_records_total{name="lua.3",url="value1"} 0 1589716338403
fluentbit_filter_add_records_total{name="rewrite_tag.2",url="value2"} 0 1589716338417`
)

func Test_Merger(t *testing.T) {

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	http.HandleFunc("/data_1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, data1)
	})
	http.HandleFunc("/data_2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, data2)
	})
	go http.Serve(listener, nil)

	urlPrefix := fmt.Sprintf("http://127.0.0.1:%d", listener.Addr().(*net.TCPAddr).Port)

	url, value1, value2 := "url", "value1", "value2"
	label1 := []*prom.LabelPair{{Name: &url, Value: &value1}}
	label2 := []*prom.LabelPair{{Name: &url, Value: &value2}}

	m := merger.New(time.Second)
	m.AddSource(urlPrefix+"/data_1", label1)
	m.AddSource(urlPrefix+"/data_2", label2)

	out := bytes.NewBuffer(make([]byte, 0))
	m.Merge(out)

	outStrs := strings.Split(out.String(), "\n")
	sort.Strings(outStrs)
	if result != strings.Join(outStrs, "\n") {
		t.Fatalf("out:\n%s\n", strings.Join(outStrs, "\n"))
	}
}
