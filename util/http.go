package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

func EncodeJson(data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func DecodeReader(r io.Reader) (map[string]interface{}, error) {
	var m map[string]interface{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func DecodeJson(data []byte) (map[string]interface{}, error) {
	var m map[string]interface{}

	err := json.Unmarshal(data, &m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func Request(method, baseUrl, path string, body io.Reader) ([]byte, int, error) {
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(25 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*20)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	url := fmt.Sprintf("%s%s", baseUrl, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	return data, resp.StatusCode, nil
}

// TruncateID returns a shorthand version of a string identifier for convenience.
// A collision with other shorthands is very unlikely, but possible.
// In case of a collision a lookup with TruncIndex.Get() will fail, and the caller
// will need to use a langer prefix, or the full-length Id.
func truncateID(id string) string {
	shortLen := 12
	if len(id) < shortLen {
		shortLen = len(id)
	}
	return id[:shortLen]
}

// GenerateRandomID returns an unique id
func GenerateRandomID() string {
	for {
		id := make([]byte, 16)
		if _, err := io.ReadFull(rand.Reader, id); err != nil {
			panic(err) // This shouldn't happen
		}
		value := hex.EncodeToString(id)
		// if we try to parse the truncated for as an int and we don't have
		// an error then the value is all numberic and causes issues when
		// used as a hostname. ref #3869
		if _, err := strconv.ParseInt(truncateID(value), 10, 32); err == nil {
			continue
		}
		return value
	}
}

func Response(data []byte, statusCode int, resp http.ResponseWriter) {
	resp.WriteHeader(statusCode)
	resp.Write(data)
}

func HandlerError(resp http.ResponseWriter, data string, httpCode int, stateCode int) {
	resp.WriteHeader(httpCode)
	data = fmt.Sprintf("{\"state\" : %d, \"error\" : \"%s\"}", stateCode, data)
	fmt.Fprintf(resp, data)
}

func AesEncrypt(encodeStr string, key []byte) (string, error) {
	encodeBytes := []byte(encodeStr)
	block, err := aes.NewCipher(key)
	if err != err {
		return "", err
	}
	blockSize := block.BlockSize()
	encodeBytes = PKCS5Padding(encodeBytes, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypted, encodeBytes)

	return base64.StdEncoding.EncodeToString(crypted), nil
}

func AesDecrypt(decodeStr string, key []byte) ([]byte, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(decodeBytes))
	blockMode.CryptBlocks(origData, decodeBytes)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func CheckValue(v interface{}, vType string) error {
	if v == nil {
		return fmt.Errorf("value is empty")
	}

	switch v.(type) {
	case int:
		if vType != "int" {
			return fmt.Errorf("%v is not %s type", v, vType)
		}
	case float64:
		if vType != "float64" {
			return fmt.Errorf("%v is not %s type", v, vType)
		}
	case string:
		if vType != "string" {
			return fmt.Errorf("%v is not %s type", v, vType)
		}
	case bool:
		if vType != "bool" {
			return fmt.Errorf("%v is not %s type", v, vType)
		}
	case []interface{}:
		if vType != "[]interface{}" {
			return fmt.Errorf("%v is not %s type", v, vType)
		}
	case map[string]interface{}:
		if vType != "map[string]interface{}" {
			return fmt.Errorf("%v is not %s type", v, vType)
		}
	default:
		return fmt.Errorf("<unknow data type>")
	}

	return nil
}

func Call(method, baseUrl, path string, body io.Reader, headers map[string][]string) ([]byte, int, error) {
	client := &http.Client{}
	fmt.Println(baseUrl + path)
	req, err := http.NewRequest(method, baseUrl+path, body)
	if err != nil {
		return nil, 408, err
	}

	req.Header.Set("User-Agent", "XENIUMD-AGENT")
	if method == "POST" {
		req.Header.Set("Content-Type", "plain/text")
	}

	if headers != nil {
		for k, v := range headers {
			req.Header[k] = v
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	dataBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	return dataBody, resp.StatusCode, nil
}
