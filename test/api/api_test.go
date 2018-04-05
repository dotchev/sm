package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dotchev/sm/rest"
	"github.com/gavv/httpexpect"
	. "github.com/onsi/ginkgo"
)

type O map[string]interface{}
type A []interface{}

func TestSM(t *testing.T) {
	os.Chdir("../..") // go back to project root
	RunSpecs(t, "Sm Suite")
}

var _ = Describe("REST API", func() {
	var server *httptest.Server
	var e *httpexpect.Expect

	BeforeSuite(func() {
		handler := rest.SMHandler()
		server = httptest.NewServer(handler)
		e = httpexpect.New(GinkgoT(), server.URL)
	})

	AfterSuite(func() {
		if server != nil {
			server.Close()
		}
	})

	Describe("Platforms", func() {
		BeforeEach(func() {
			By("remove all platforms")

			resp := e.GET("/v1/platforms").
				Expect().Status(http.StatusOK).JSON().Object()
			for _, val := range resp.Value("platforms").Array().Iter() {
				id := val.Object().Value("id").String().Raw()
				e.DELETE("/v1/platforms/" + id).
					Expect().Status(http.StatusOK)
			}
		})

		Describe("GET", func() {
			It("returns 404 if platform does not exist", func() {
				e.GET("/v1/platforms/999").
					Expect().
					Status(http.StatusNotFound).
					JSON().Object().Keys().Contains("error", "description")
			})

			It("returns the platform with given id", func() {
				platform := O{
					"name":        "cf-05",
					"type":        "cf",
					"description": "CF somewhere",
				}
				reply := e.POST("/v1/platforms").WithJSON(platform).
					Expect().Status(http.StatusCreated).JSON().Object()
				id := reply.Value("id").String().Raw()
				platform["id"] = id

				e.GET("/v1/platforms/" + id).
					Expect().
					Status(http.StatusOK).
					JSON().Equal(platform)
			})
		})

		Describe("GET All", func() {
			It("returns empty array if no platforms exist", func() {
				e.GET("/v1/platforms").
					Expect().
					Status(http.StatusOK).
					JSON().Object().Value("platforms").Array().Empty()
			})

			It("returns all the platforms", func() {
				platforms := A{}

				addPlatform := func(name string, ptype string, description string) {
					platform := O{
						"name":        name,
						"type":        ptype,
						"description": description,
					}
					reply := e.POST("/v1/platforms").WithJSON(platform).
						Expect().Status(http.StatusCreated).JSON().Object()
					id := reply.Value("id").String().Raw()
					platform["id"] = id
					platforms = append(platforms, platform)

					e.GET("/v1/platforms").
						Expect().
						Status(http.StatusOK).
						JSON().Object().Value("platforms").Array().ContainsOnly(platforms...)
				}
				addPlatform("cf-05", "cf", "CF in EU")
				addPlatform("kube-06", "k8s", "Kubernetes in US")
			})
		})

		Describe("POST", func() {
			It("returns 400 if input is not valid JSON", func() {
				e.POST("/v1/platforms").
					WithText("text").
					Expect().Status(http.StatusBadRequest)
				e.POST("/v1/platforms").
					WithText("invalid json").
					WithHeader("content-type", "application/json").
					Expect().Status(http.StatusBadRequest)
			})

			It("returns 400 if mandatory field is missing", func() {
				e.POST("/v1/platforms").
					WithJSON(O{"name": "cf-05"}).
					Expect().Status(http.StatusBadRequest)
				e.POST("/v1/platforms").
					WithJSON(O{"type": "cf"}).
					Expect().Status(http.StatusBadRequest)
			})

			It("succeeds if optional fields are skipped", func() {
				platform := O{
					"name": "cf-05",
					"type": "cf",
				}

				reply := e.POST("/v1/platforms").
					WithJSON(platform).
					Expect().Status(http.StatusCreated).JSON().Object()

				platform["id"] = reply.Value("id").String().Raw()
				platform["description"] = ""
				reply.Equal(platform)
			})

			It("generates a valid platform id", func() {
				platform := O{
					"name":        "cf-05",
					"type":        "cf",
					"description": "CF in EU",
				}

				reply := e.POST("/v1/platforms").
					WithJSON(platform).
					Expect().Status(http.StatusCreated).JSON().Object()

				id := reply.Value("id").String().NotEmpty().Raw()
				platform["id"] = id
				reply.Equal(platform)
				e.GET("/v1/platforms/" + id).
					Expect().Status(http.StatusOK).JSON().Equal(platform)
			})

			It("uses provided platform id", func() {
				platform := O{
					"id":          "abc-123",
					"name":        "cf-05",
					"type":        "cf",
					"description": "CF in EU",
				}

				reply := e.POST("/v1/platforms").
					WithJSON(platform).
					Expect().Status(http.StatusCreated).JSON().Object()

				reply.Equal(platform)
				e.GET("/v1/platforms/abc-123").
					Expect().Status(http.StatusOK).JSON().Equal(platform)
			})

			It("returns 400 if duplicate id is provided", func() {
				e.POST("/v1/platforms").
					WithJSON(O{
						"id":   "abc-123",
						"name": "cf-05",
						"type": "cf",
					}).
					Expect().Status(http.StatusCreated)

				e.POST("/v1/platforms").
					WithJSON(O{
						"id":   "abc-123",
						"name": "cf-06",
						"type": "cff",
					}).
					Expect().Status(http.StatusBadRequest)
			})

			It("returns 400 if duplicate name is provided", func() {
				e.POST("/v1/platforms").
					WithJSON(O{
						"name": "cf-05",
						"type": "cf",
					}).
					Expect().Status(http.StatusCreated)

				e.POST("/v1/platforms").
					WithJSON(O{
						"name": "cf-05",
						"type": "cff",
					}).
					Expect().Status(http.StatusBadRequest)
			})
		})

		Describe("PATCH", func() {
			It("returns 400 if input is not valid JSON", func() {
				e.PATCH("/v1/platforms/1").
					WithText("text").
					Expect().Status(http.StatusBadRequest)
				e.PATCH("/v1/platforms/1").
					WithText("invalid json").
					WithHeader("content-type", "application/json").
					Expect().Status(http.StatusBadRequest)
			})

			It("returns 404 if platform does not exist", func() {
				e.PATCH("/v1/platforms/999").
					WithJSON(O{}).
					Expect().
					Status(http.StatusNotFound)
			})

			It("updates existing platform", func() {
				By("create initial platform")

				platform := O{
					"id":          "1",
					"name":        "cf-05",
					"type":        "cf",
					"description": "CF in EU",
				}
				e.POST("/v1/platforms").
					WithJSON(platform).
					Expect().Status(http.StatusCreated)

				By("all fields are optional")

				e.PATCH("/v1/platforms/1").
					WithJSON(O{}).
					Expect().
					Status(http.StatusOK).JSON().Equal(platform)

				By("ignores given id")

				e.PATCH("/v1/platforms/1").
					WithJSON(O{"id": "2"}).
					Expect().
					Status(http.StatusOK).JSON().Equal(platform)

				By("updates only provided fields")

				platform["name"] = "cf6"
				e.PATCH("/v1/platforms/1").
					WithJSON(O{"name": "cf6"}).
					Expect().
					Status(http.StatusOK).JSON().Equal(platform)

				By("GET returns the same state")

				e.GET("/v1/platforms/1").
					Expect().Status(http.StatusOK).JSON().Equal(platform)
			})
		})

		Describe("DELETE", func() {
			It("deletes given platform", func() {
				By("create initial platform")

				platform := O{
					"id":   "1",
					"name": "cf-05",
					"type": "cf",
				}
				e.POST("/v1/platforms").
					WithJSON(platform).
					Expect().Status(http.StatusCreated)

				By("delete the platform")

				e.DELETE("/v1/platforms/1").
					Expect().Status(http.StatusOK).JSON().Equal(O{})

				By("the platform is gone")

				e.GET("/v1/platforms/1").Expect().Status(http.StatusNotFound)
				e.GET("/v1/platforms").Expect().Status(http.StatusOK).
					JSON().Object().Value("platforms").Array().Empty()
				e.DELETE("/v1/platforms/1").Expect().Status(http.StatusNotFound)
			})
		})
	})
})
