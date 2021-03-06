package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FirebaseManifest struct {
	Public    string `json:"public"`
	Redirects []struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		Type        int    `json:"type,omitempty"`
	} `json:"redirects"`
	Rewrites []struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	} `json:"rewrites"`
	Headers []struct {
		Source  string `json:"source"`
		Headers []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"headers"`
	} `json:"headers"`
	Hosting *FirebaseManifest `json:"Hosting"`
}

func (mf FirebaseManifest) processRedirects(w http.ResponseWriter, r *http.Request) (bool, error) {
	for _, redirect := range mf.Redirects {
		pattern, err := CompileExtGlob(redirect.Source)
		if err != nil {
			return false, fmt.Errorf("Invalid redirect extglob %s: %s", redirect.Source, err)
		}
		if pattern.MatchString(r.URL.Path) {
			http.Redirect(w, r, redirect.Destination, redirect.Type)
			return true, nil
		}
	}
	if mf.Hosting != nil {
		return mf.Hosting.processRedirects(w, r)
	}
	return false, nil
}

func (mf FirebaseManifest) processRewrites(r *http.Request) error {
	for _, rewrite := range mf.Rewrites {
		pattern, err := CompileExtGlob(rewrite.Source)
		if err != nil {
			return fmt.Errorf("Invalid rewrite extglob %s: %s", rewrite.Source, err)
		}
		if pattern.MatchString(r.URL.Path) {
			r.URL.Path = strings.TrimSuffix(rewrite.Destination, "index.html")
			return nil
		}
	}
	if mf.Hosting != nil {
		return mf.Hosting.processRewrites(r)
	}
	return nil
}

func (mf FirebaseManifest) processHosting(w http.ResponseWriter, r *http.Request) error {
	for _, headerSet := range mf.Headers {
		pattern, err := CompileExtGlob(headerSet.Source)
		if err != nil {
			return fmt.Errorf("Invalid hosting.header extglob %s: %s", headerSet.Source, err)
		}
		if pattern.MatchString(r.URL.Path) {
			for _, header := range headerSet.Headers {
				w.Header().Set(header.Key, header.Value)
			}
		}
	}
	if mf.Hosting != nil {
		return mf.Hosting.processHosting(w, r)
	}
	return nil
}

func processWithConfig(w http.ResponseWriter, r *http.Request, config string) string {
	dir := "."
	mf, err := readManifest(config)
	if err != nil {
		log.Printf("Could read Firebase file %s: %s", config, err)
		return dir
	}
	if mf.Public != "" {
		dir = mf.Public
	}
	if mf.Hosting != nil && mf.Hosting.Public != "" {
		dir = mf.Hosting.Public
	}

	done, err := mf.processRedirects(w, r)
	if err != nil {
		log.Printf("Processing redirects failed: %s", err)
		return dir
	}
	if done {
		return dir
	}

	// Rewrites only happen if the target file does not exist
	if _, err = os.Stat(filepath.Join(dir, r.URL.Path)); err != nil {
		err = mf.processRewrites(r)
		if err != nil {
			log.Printf("Processing rewrites failed: %s", err)
			return dir
		}
	}

	err = mf.processHosting(w, r)
	if err != nil {
		log.Printf("Processing rewrites failed: %s", err)
		return dir
	}

	return dir
}

func readManifest(path string) (FirebaseManifest, error) {
	fmf := FirebaseManifest{}
	f, err := os.Open(path)
	if err != nil {
		return fmf, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&fmf)
	return fmf, err
}
