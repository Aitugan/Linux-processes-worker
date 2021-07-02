package cmd

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func SendRequest(method, url, certsDir string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, body)
	trace := &httptrace.ClientTrace{
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	log.Info("reading client keypair")
	cert, err := tls.LoadX509KeyPair(certsDir+"/client-cert.pem", certsDir+"/client-key.pem")
	if err != nil {
		log.Fatalf("client: loadkeys: %s", err)
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	trustStore, err := ioutil.ReadFile(certsDir + "/CA-cert.pem")
	if err != nil {
		log.Fatalf("failed to get client trust store %s", err)
		return nil, err
	}

	caCertPool.AppendCertsFromPEM(trustStore)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: false,
		RootCAs:            caCertPool,
	}

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// b, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	return res, nil
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A start command",
	Long:  "A start command",
	Run: func(cmd *cobra.Command, args []string) {
		certsDir, _ := cmd.Flags().GetString("dir")
		fmt.Println(certsDir)
		if len(args) == 0 {
			fmt.Println("An executable command is required")
		}
		values := make(map[string][]string)
		fmt.Println(args[0])
		values["command"] = args
		json_data, err := json.Marshal(values)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := SendRequest("POST", "https://localhost:8080/work/start", certsDir, bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A start command",
	Long:  "A start command",
	Run: func(cmd *cobra.Command, args []string) {
		certsDir, _ := cmd.Flags().GetString("dir")
		if len(args) == 0 {
			fmt.Println("A command id is required")
		}
		values := make(map[string]string)
		values["commandID"] = args[0]
		json_data, err := json.Marshal(values)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := SendRequest("PUT", "https://localhost:8080/work/stop/"+args[0], certsDir, bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	},
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "A start command",
	Long:  "A start command",
	Run: func(cmd *cobra.Command, args []string) {
		certsDir, _ := cmd.Flags().GetString("dir")
		fmt.Println(certsDir)
		if len(args) == 0 {
			fmt.Println("An executable command is required")
		}
		values := make(map[string]string)
		fmt.Println(args[0])
		values["commandID"] = args[0]
		json_data, err := json.Marshal(values)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := SendRequest("POST", "https://localhost:8080/work/log", certsDir, bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	},
}

var queryStatusCmd = &cobra.Command{
	Use:   "queryStatus",
	Short: "A start command",
	Long:  "A start command",
	Run: func(cmd *cobra.Command, args []string) {
		certsDir, _ := cmd.Flags().GetString("dir")
		fmt.Println(certsDir)
		if len(args) == 0 {
			fmt.Println("An executable command is required")
		}
		// values := make(map[string]string)
		fmt.Println(args[0])

		resp, err := SendRequest("GET", "https://localhost:8080/work/query-status/"+args[0], certsDir, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	},
}
