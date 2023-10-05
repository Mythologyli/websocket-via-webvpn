package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"runtime/debug"
	"strings"
)

func login(username string, password string) (string, error) {
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}

	addr := "https://webvpn.zju.edu.cn/por/login_auth.csp?apiversion=1"

	resp, err := c.Get(addr)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", err
	}

	twfId := string(regexp.MustCompile(`<TwfID>(.*)</TwfID>`).FindSubmatch(buf.Bytes())[1])
	rsaKey := string(regexp.MustCompile(`<RSA_ENCRYPT_KEY>(.*)</RSA_ENCRYPT_KEY>`).FindSubmatch(buf.Bytes())[1])

	csrfCode := string(regexp.MustCompile(`<CSRF_RAND_CODE>(.*)</CSRF_RAND_CODE>`).FindSubmatch(buf.Bytes())[1])
	password += "_" + csrfCode

	pubKey := rsa.PublicKey{}
	pubKey.E = 65537
	modulus := big.Int{}
	modulus.SetString(rsaKey, 16)
	pubKey.N = &modulus

	encryptedPassword, err := rsa.EncryptPKCS1v15(rand.Reader, &pubKey, []byte(password))
	if err != nil {
		return "", err
	}
	encryptedPasswordHex := hex.EncodeToString(encryptedPassword)

	addr = "https://webvpn.zju.edu.cn/por/login_psw.csp?anti_replay=1&encrypt=1&apiversion=1"

	form := url.Values{
		"svpn_rand_code":    {""},
		"mitm_result":       {""},
		"svpn_req_randcode": {csrfCode},
		"svpn_name":         {username},
		"svpn_password":     {encryptedPasswordHex},
	}

	req, err := http.NewRequest("POST", addr, strings.NewReader(form.Encode()))
	req.Header.Set("Cookie", "TWFID="+twfId)

	resp, err = c.Do(req)
	if err != nil {
		debug.PrintStack()
		return "", err
	}

	buf.Reset()
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if !strings.Contains(string(buf.Bytes()), "radius auth succ") {
		return "", errors.New("login failed")
	}

	return twfId, nil
}
