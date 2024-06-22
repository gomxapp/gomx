package internal

import (
	"fmt"
	"github.com/gomxapp/gomx/pkg/router"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var tree *RouteTree

func ExpectEqual[T comparable](t *testing.T, actual T, expected T) {
	if actual != expected {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			t.Log("Error getting caller info")
		} else {
			t.Logf("%s:%d", filepath.Base(file), line)
		}
		t.Errorf("\nActual: %v\nExpected: %v", actual, expected)
	}
}

func makeMockNode(name string) *RouteTree {
	return createNode(name, router.GET,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			log.Println("Handler for " + name)
		}), nil)
}

func makeMockWildNode(name string) *RouteTree {
	return createNode(name, router.GET,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Handler for " + name)
			wildData := r.PathValue(strings.Trim(name, "{}"))
			log.Printf("Value for %s: %s", name, wildData)
			_, _ = w.Write([]byte(wildData))
		}), nil)
}

func makeMockNotFoundHandlerNode(name string) *RouteTree {
	return createNode(name, router.GET,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			log.Println("Handler for " + name)
		}), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Route at " + r.URL.Path + " not found. Using notFoundHandler from " + name)
			_, _ = w.Write([]byte(name))
		}),
	)
}

func makeTestTree() *RouteTree {
	tree := createRoot()
	getRootNode := makeMockNode("")
	_ = tree.AddChild(getRootNode)
	aNode := makeMockNotFoundHandlerNode("a")
	_ = getRootNode.AddChild(aNode)
	bNode := makeMockNode("b")
	_ = aNode.AddChild(bNode)
	cNode := makeMockNotFoundHandlerNode("c")
	_ = bNode.AddChild(cNode)
	dNode := makeMockNode("d")
	_ = cNode.AddChild(dNode)
	eNode := makeMockNode("e")
	_ = getRootNode.AddChild(eNode)
	wildNode := makeMockWildNode("{z}")
	_ = eNode.AddChild(wildNode)
	fNode := makeMockNode("f")
	_ = wildNode.AddChild(fNode)
	gNode := makeMockNode("g")
	_ = getRootNode.AddChild(gNode)
	return tree
}

func TestMain(m *testing.M) {
	tree = makeTestTree()
	fmt.Println("Running tests with tree:")
	fmt.Println(tree)
	m.Run()
}

func TestRouteTree(t *testing.T) {
	var incomingRequestPath string
	var incomingRequestMethod router.Method
	var mockRequest *http.Request
	var mockWriter *httptest.ResponseRecorder

	// ------------- TEST BASIC ROUTING
	t.Run("test basic routing", func(t *testing.T) {
		incomingRequestPath = "/a/b/c"
		incomingRequestMethod = router.GET
		mockRequest = httptest.NewRequest(http.MethodGet, incomingRequestPath, nil)
		t.Log("Testing path: " + incomingRequestPath)

		closestNode, matchLevel := tree.FindClosestMatchingNode(incomingRequestPath, incomingRequestMethod)
		closestNodePath := closestNode.GetPath()
		expectedMatchLevel := ExactMatch
		expectedPath := "/a/b/c/"
		ExpectEqual(t, matchLevel, expectedMatchLevel)
		ExpectEqual(t, closestNodePath, expectedPath)
		closestNode.ServeHTTP(nil, mockRequest)

		incomingRequestPath = "/"
		incomingRequestMethod = router.GET
		mockRequest = httptest.NewRequest(http.MethodGet, incomingRequestPath, nil)
		t.Log("Testing path: " + incomingRequestPath)

		closestNode, matchLevel = tree.FindClosestMatchingNode(incomingRequestPath, incomingRequestMethod)
		closestNodePath = closestNode.GetPath()
		expectedMatchLevel = ExactMatch
		expectedPath = "/"
		ExpectEqual(t, matchLevel, expectedMatchLevel)
		ExpectEqual(t, closestNodePath, expectedPath)
		closestNode.ServeHTTP(nil, mockRequest)
	})

	// ------------- TEST WILDCARD ROUTES
	t.Run("test wildcard routes", func(t *testing.T) {
		wildData := "test"
		incomingRequestPath = "/e/" + wildData
		incomingRequestMethod = router.GET
		mockRequest = httptest.NewRequest(http.MethodGet, incomingRequestPath, nil)
		t.Log("Testing path: " + incomingRequestPath)

		closestNode, matchLevel := tree.FindClosestMatchingNode(incomingRequestPath, incomingRequestMethod)
		closestNodePath := closestNode.GetPath()
		expectedMatchLevel := WildMatch
		expectedPath := "/e/z/"
		ExpectEqual(t, matchLevel, expectedMatchLevel)
		ExpectEqual(t, closestNodePath, expectedPath)
		ExpectEqual(t, closestNode.wildData, wildData)
	})

	// ------------- TEST NOT FOUND HANDLING
	t.Run("test not found handling", func(t *testing.T) {
		incomingRequestPath = "/a/b/FAKEPATH"
		incomingRequestMethod = router.GET
		mockWriter = httptest.NewRecorder()
		mockRequest = httptest.NewRequest(string(incomingRequestMethod), incomingRequestPath, nil)
		t.Log("Testing path: " + incomingRequestPath)
		closestNode, matchLevel := tree.FindClosestMatchingNode(incomingRequestPath, incomingRequestMethod)
		expectedMatchLevel := NoMatch
		ExpectEqual(t, matchLevel, expectedMatchLevel)
		closestNotFoundHandler := closestNode.FindClosestNotFoundHandler()
		closestNodePath := closestNotFoundHandler.GetPath()
		expectedPath := "/a/"
		ExpectEqual(t, closestNodePath, expectedPath)
		closestNotFoundHandler.notFoundHandler.ServeHTTP(mockWriter, mockRequest)
		ExpectEqual(t, mockWriter.Body.String(), "a")

		incomingRequestPath = "/a/b/c/d/asdf/asdf"
		incomingRequestMethod = router.GET
		mockWriter = httptest.NewRecorder()
		mockRequest = httptest.NewRequest(string(incomingRequestMethod), incomingRequestPath, nil)
		t.Log("Testing path: " + incomingRequestPath)
		closestNode, matchLevel = tree.FindClosestMatchingNode(incomingRequestPath, incomingRequestMethod)
		expectedMatchLevel = NoMatch
		ExpectEqual(t, matchLevel, expectedMatchLevel)
		closestNotFoundHandler = closestNode.FindClosestNotFoundHandler()
		closestNodePath = closestNotFoundHandler.GetPath()
		expectedPath = "/a/b/c/"
		ExpectEqual(t, closestNodePath, expectedPath)
		closestNotFoundHandler.notFoundHandler.ServeHTTP(mockWriter, mockRequest)
		ExpectEqual(t, mockWriter.Body.String(), "c")
	})
}
