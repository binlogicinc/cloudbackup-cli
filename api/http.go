package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const ISO_8601_FORMAT = "2006-01-02T15:04:05-0700"

type signedHTTPClient struct {
	http.Client
	AccessKey string
	SecretKey string
}

func NewSignedHTTPClient(accessKey, secretKey string, timeoutSecs int) *signedHTTPClient {
	return &signedHTTPClient{
		Client: http.Client{
			Timeout: time.Second * time.Duration(timeoutSecs),
		},
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

func (cli *signedHTTPClient) SignedGet(url string, headers map[string]string) (*http.Response, error) {
	return cli.SignedDo("GET", url, nil, headers)
}

func (cli *signedHTTPClient) SignedDelete(url string, headers map[string]string) (*http.Response, error) {
	return cli.SignedDo("DELETE", url, nil, headers)
}

func (cli *signedHTTPClient) SignedPost(url string, body io.Reader,
	headers map[string]string) (*http.Response, error) {

	return cli.SignedDo("POST", url, body, headers)
}

func (cli *signedHTTPClient) SignedDo(verb, url string, body io.Reader,
	headers map[string]string) (*http.Response, error) {

	req, err := http.NewRequest(verb, url, body)

	if err != nil {
		return nil, err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	cli.sign(req, body)

	fmt.Println(req.Method, req.URL.String(), req.Header)

	resp, err := cli.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (cli *signedHTTPClient) sign(req *http.Request, body io.Reader) error {
	var buff bytes.Buffer

	date := time.Now().UTC().Format(ISO_8601_FORMAT)
	req.Header.Add("date", date)
	req.Header.Add("bl-access-key", cli.AccessKey)

	buff.WriteString(req.Method)
	buff.WriteString("\n")
	buff.WriteString(req.URL.String())
	buff.WriteString("\n")
	buff.WriteString(date)
	buff.WriteString("\n")
	buff.WriteString(cli.AccessKey)
	buff.WriteString("\n")

	if body != nil {
		bodyBytes, err := ioutil.ReadAll(body)

		if err != nil {
			return err
		}

		sum := md5.Sum(bodyBytes)
		buff.WriteString(hex.EncodeToString(sum[:]))

		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	hash := hmac.New(sha256.New, []byte(cli.SecretKey))
	_, err := hash.Write(buff.Bytes())

	if err != nil {
		return err
	}

	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	req.Header.Add("Authorization", "BL "+signature)
	req.Header.Add("bl-msg", strings.Replace(buff.String(), "\n", "\\n", -1))

	return nil
}

func (cli *signedHTTPClient) postJSON(url string, i interface{}) (val map[string]interface{}, err error) {
	b, err := json.Marshal(i)

	if err != nil {
		err = wrap("while marshalling interface", err)
		return
	}

	resp, err := cli.SignedPost(url, bytes.NewBuffer(b), defaultHeaders)

	if err != nil {
		err = wrap("while sending signed post", err)
		return
	}

	body, val, err := cli.parseResponseJSON(resp)

	if err != nil {
		err = wrap("while parsing json response", err)
		return
	}

	_, err = cli.isJSONResponseOk(body, val)

	return
}

func (cli *signedHTTPClient) parseResponseJSON(resp *http.Response) (body []byte, val map[string]interface{}, err error) {
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		err = wrap("while reading body", err)
		return
	}

	err = json.Unmarshal(body, &val)

	if err != nil {
		err = wrap("while unmarshalling body "+string(body), err)
		return
	}

	return
}

func (cli *signedHTTPClient) isJSONResponseOk(body []byte, val map[string]interface{}) (bool, error) {
	if status, ok := val["status"]; ok { //if status is present
		if s, ok2 := status.(string); s != "ok" {
			if !ok2 { //if status couldnt be parsed as string
				return false, fmt.Errorf("Unexpected response %s", string(body))
			} else { //status was a string but it was not ok
				return false, fmt.Errorf("%s", val["message"])
			}
		} else { //status is equal to ok, so no error
			return true, nil
		}
	} else { //status key is not present
		return false, fmt.Errorf("Unexpected response %s", string(body))
	}
}
