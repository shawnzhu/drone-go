package drone

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"github.com/drone/drone/shared/sshutil"
	"encoding/json"
	"testing"

	"github.com/franela/goblin"
)

func Test_Client(t *testing.T) {
	mockToken := "fake_token"
	g := goblin.Goblin(t)

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()
        g.Assert(qs.Get("access_token")).Equal(mockToken)

		defer r.Body.Close()
		switch r.URL.Path {
		case "/api/repos/host/owner/name":
			switch r.Method {
			case "PUT":
				in := struct {
					PostCommit  *bool   `json:"post_commits"`
					PullRequest *bool   `json:"pull_requests"`
					Privileged  *bool   `json:"privileged"`
					Params      *string `json:"params"`
					Timeout     *int64  `json:"timeout"`
					PublicKey   *string `json:"public_key"`
					PrivateKey  *string `json:"private_key"`
				}{}
				if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				//fmt.Printf(*in.PublicKey)
				//fmt.Printf(*in.PrivateKey)

				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"fake json string"}`)
			}
		}
    }))
    defer ts.Close()

	client := NewClient(mockToken, ts.URL)
/*
    g.Describe("Drone Client", func() {
        g.It("Should send access_token over query string", func() {
            resp, err := client.do("GET", "/foo")
            g.Assert(err).Equal(nil)
			g.Assert(resp.Header.Get("Content-Type")).Equal("application/json")
        })
    })
*/
	g.Describe("Drone repos service", func() {
		repos := client.Repos

		g.It("capable to setup keypair", func() {
			privateKey, err := sshutil.GeneratePrivateKey()
			g.Assert(err).Equal(nil)

			pub  := sshutil.MarshalPublicKey(&privateKey.PublicKey)
			priv := sshutil.MarshalPrivateKey(privateKey)

			//fmt.Printf(pub)
			//fmt.Printf(priv)

			err = repos.SetKey("host", "owner", "name", pub, priv)
			g.Assert(err).Equal(nil)
		})
	})
}
