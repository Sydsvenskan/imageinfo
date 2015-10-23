package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"

	"github.com/satori/go.uuid"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	help := flag.Bool("help", false, "Show this help")
	file := flag.String("file", "", "A file to calculate the uuid for")
	base := flag.String("base", "cmbr://image/", "The base URI for the entity")
	port := flag.Int("port", 0, "A port to listen for image requests on")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	baseURL, err := url.Parse(*base)
	if err != nil {
		log.Fatalf("Could not parse the base URI %q: %s", *base, err)
	}

	var reader io.Reader

	if *file != "" {
		fileReader, err := os.Open(*file)
		if err != nil {
			log.Fatalf("Could not open the specified file %q for reading: %s", *file, err)
		}
		defer fileReader.Close()
		reader = fileReader
	}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		reader = os.Stdin
	}

	if reader != nil {
		enc := json.NewEncoder(os.Stdout)
		res, err := getFileInfo(baseURL, reader)
		if err != nil {
			log.Fatalf("Failed to get file info: %s", err)
		}

		err = enc.Encode(res)
		if err != nil {
			log.Fatalf("Failed to encode result: %s", err)
		}
	}

	if *port != 0 {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			enc := json.NewEncoder(w)
			w.Header().Add("Content-type", "application/json")

			results := make(map[string]interface{})

			reader, err := r.MultipartReader()
			if err != nil {
				badRequest(w, err.Error())
				return
			}
			for {
				part, err := reader.NextPart()

				if err == io.EOF {
					break
				}
				if err != nil {
					badRequest(w, err.Error())
					return
				}

				if len(part.Header["Content-Disposition"]) == 0 {
					badRequest(w, "Missing Content-Disposition header in multipart part")
					return
				}
				_, params, err := mime.ParseMediaType(part.Header["Content-Disposition"][0])
				if err != nil {
					badRequest(w, "Failed to parse Content-Disposition %q", part.Header["Content-Disposition"][0])
					return
				}

				fieldName, ok := params["name"]
				if !ok {
					badRequest(w, "Missing Content-Disposition name in %q", part.Header["Content-Disposition"][0])
					return
				}

				res, err := getFileInfo(baseURL, part)
				if err != nil {
					results[fieldName] = map[string]interface{}{
						"error": err.Error(),
					}
					log.Printf("Failed to get file info: %s", err)
				} else {
					results[fieldName] = res
				}
			}

			w.WriteHeader(http.StatusOK)
			enc.Encode(results)
		})

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	}
}

type result struct {
	URI    string    `json:"uri"`
	UUID   uuid.UUID `json:"uuid"`
	Width  int       `json:"width"`
	Height int       `json:"height"`
}

func badRequest(w http.ResponseWriter, message string, values ...interface{}) {
	enc := json.NewEncoder(w)

	w.WriteHeader(http.StatusBadRequest)
	enc.Encode(map[string]string{
		"error": fmt.Sprintf(message, values...),
	})
}

func getFileInfo(base *url.URL, reader io.Reader) (*result, error) {
	h := sha1.New()
	buf := make([]byte, 1024)

	tee := io.TeeReader(reader, h)
	res := result{}

	config, format, err := image.DecodeConfig(tee)
	if err != nil {
		return nil, err
	}
	res.Width = config.Width
	res.Height = config.Height

	for ; err != io.EOF; _, err = tee.Read(buf) {
		if err != nil {
			return nil, err
		}
	}

	hash := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	hashURL, _ := url.Parse(hash + "." + format)
	res.URI = base.ResolveReference(hashURL).String()
	res.UUID = uuid.NewV5(uuid.NamespaceURL, res.URI)
	return &res, nil
}
