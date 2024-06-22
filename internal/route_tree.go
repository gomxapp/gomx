package internal

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type MatchLevel int

const (
	NoMatch          MatchLevel = iota // target: /a/b/c tree: /a/d/e
	PartialWildMatch                   // target: /a/b/c tree: /a/{z}
	PartialMatch                       // target: /a/b/c tree: /a/b
	WildMatch                          // target: /a/b/c tree: /a/{z}/c
	ExactMatch                         // target: /a/b/c tree: /a/b/c
)

// ----------------- ROUTE TREE WRAPPER

type RouteTreeWrapper struct {
	Tree        *RouteTree
	closestNode *RouteTree
	matchLvl    MatchLevel
}

func (wrapper *RouteTreeWrapper) ServeNotFound(w http.ResponseWriter, r *http.Request) {
	if n := wrapper.closestNode; n != nil && n.notFoundHandler != nil {
		n.notFoundHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (wrapper *RouteTreeWrapper) setClosestMatch(path string, method string) {
	wrapper.closestNode, wrapper.matchLvl = wrapper.Tree.FindClosestMatchingNode(path, method)
}

func (wrapper *RouteTreeWrapper) setPathValues(r *http.Request) {
	if wrapper.matchLvl != WildMatch {
		return
	}
	match := wrapper.closestNode
	nodes, err := match.GetPathFromRoot(false)
	if err != nil {
		panic("error getting path from root")
	}
	for _, n := range nodes {
		if n.isWild {
			r.SetPathValue(n.pathPart, n.wildData)
		}
	}
}

func (wrapper *RouteTreeWrapper) ContainsExactMatch(request *http.Request) bool {
	path := request.URL.EscapedPath()
	method := request.Method
	wrapper.setClosestMatch(path, method)
	wrapper.setPathValues(request)
	return wrapper.matchLvl >= WildMatch && wrapper.closestNode != nil && wrapper.closestNode.handler != nil
}

func (wrapper *RouteTreeWrapper) ServeClosestMatch(w http.ResponseWriter, r *http.Request) {
	wrapper.closestNode.ServeHTTP(w, r)
}

func (wrapper *RouteTreeWrapper) String() string {
	return wrapper.Tree.String()
}

// ----------------- ROUTE TREE

type RouteTree struct {
	pathPart        string
	method          string
	handler         http.Handler
	notFoundHandler http.Handler
	parent          *RouteTree
	children        []*RouteTree
	isWild          bool
	wildData        string
}

func createRoot() *RouteTree {
	return createNode("", "", nil, nil)
}

func createNode(path string, method string, handler http.Handler, errorHandler http.Handler) *RouteTree {
	isWild := false
	if strings.HasPrefix(path, wildcardPrefix) &&
		strings.HasSuffix(path, wildcardSuffix) {
		path = strings.TrimPrefix(path, wildcardPrefix)
		path = strings.TrimSuffix(path, wildcardSuffix)
		isWild = true
	}
	path = strings.Trim(path, "/")
	return &RouteTree{
		pathPart:        path,
		method:          method,
		handler:         handler,
		notFoundHandler: errorHandler,
		parent:          nil,
		children:        make([]*RouteTree, 0),
		isWild:          isWild,
	}
}

func (tree *RouteTree) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if tree != nil && tree.handler != nil {
		tree.handler.ServeHTTP(w, r)
	}
}

// IsRoot returns whether the node is the root node.
func (tree *RouteTree) IsRoot() bool {
	return tree.parent == nil
}

// Depth returns the depth of the node, not including the root.
//
// "" - "/" - "home/" => depth 2
func (tree *RouteTree) Depth() int {
	if tree.IsRoot() {
		return 0
	}
	return 1 + tree.parent.Depth()
}

// Go returns the child of tree with the given path part and matching method,
// or nil if one was not found. Go takes wildcards into account.
func (tree *RouteTree) Go(nextPathPart string, method string) *RouteTree {
	for _, child := range tree.children {
		if child.method != method {
			continue
		}
		if pathPartIsWildcard(nextPathPart) && child.isWild {
			return child
		}
		if child.pathPart == nextPathPart {
			return child
		}
	}
	return nil
}

// AddChild adds a child node to this node. If a child exists with the same
// path and method, but its handlers are nil, AddChild replaces that node.
//
// Returns an error if tree or child is nil, child already has a parent,
// or a child already exists with handlers.
func (tree *RouteTree) AddChild(child *RouteTree) error {
	if tree == nil {
		return errors.New("adding to nil node")
	}
	if child == nil {
		return errors.New("adding nil child node")
	}
	if !child.IsRoot() {
		return errors.New("child already had a parent node")
	}
	if c := tree.Go(child.pathPart, child.method); c != nil {
		// child exists but can be overwritten
		if c.handler == nil && c.notFoundHandler == nil {
			c.handler = child.handler
			c.notFoundHandler = child.notFoundHandler
			return nil
		}
		// child exists but cannot be overwritten
		return errors.New(fmt.Sprintf(
			"child already exists with pathPart=%s and method=%s", child.pathPart, child.method,
		))
	}
	tree.children = append(tree.children, child)
	child.parent = tree
	return nil
}

// AddRelativeChild adds a child node to the tree given a relative path from tree.
// If contains sections missing from the tree, new nodes will be created with
// nil handlers.
//
// If successful, AddRelativeChild returns child=the child with handlers and error=nil.
// If not, child=nil.
func (tree *RouteTree) AddRelativeChild(relPath string, method string, handler http.Handler, notFoundHandler http.Handler) (*RouteTree, error) {
	if tree == nil {
		return nil, errors.New("adding to nil node")
	}
	curr := tree
	parts := strings.Split(relPath, "/")
	for _, pathPart := range parts {
		newNode := createNode(pathPart, method, nil, nil)
		err := curr.AddChild(newNode)
		if err != nil {
			// only here if child already exists
			curr = curr.Go(newNode.pathPart, newNode.method)
			continue
		}
		curr = newNode
	}
	curr.handler = handler
	curr.notFoundHandler = notFoundHandler
	return curr, nil
}

// GetPath returns the full match path at the current node. Except for the root
// node, GetPath appends a trailing-slash to the full path.
func (tree *RouteTree) GetPath() string {
	var helper func(*RouteTree) string
	helper = func(tree *RouteTree) string {
		if tree == nil {
			return ""
		}
		if tree.IsRoot() {
			return ""
		}
		suffix := ""
		if !strings.HasSuffix(tree.pathPart, "/") {
			suffix = "/"
		}
		return helper(tree.parent) + tree.pathPart + suffix
	}
	out := helper(tree)
	if !strings.HasSuffix(out, "/") {
		out = out + "/"
	}
	return out
}

// GetPathFromRoot returns all nodes to reach this node as a *RouteTree slice
// ordered starting from the root.
func (tree *RouteTree) GetPathFromRoot(includeRoot bool) ([]*RouteTree, error) {
	if tree == nil {
		return nil, errors.New("tree is nil")
	}
	if tree.IsRoot() {
		if includeRoot {
			return []*RouteTree{tree}, nil
		}
		return []*RouteTree{}, nil
	}
	p, err := tree.parent.GetPathFromRoot(includeRoot)
	if err != nil {
		return nil, err
	}
	return append(p, tree), nil
}

const (
	wildcardPrefix = "{"
	wildcardSuffix = "}"
)

func pathPartIsWildcard(pathPart string) bool {
	return strings.HasPrefix(pathPart, wildcardPrefix) &&
		strings.HasSuffix(pathPart, wildcardSuffix)
}

// compareWithPath returns how much tree.GetPath() matches targetPath.
func (tree *RouteTree) compareWithPath(targetPath string) MatchLevel {
	treePath := tree.GetPath()
	parts := strings.Split(treePath, "/")
	targetParts := strings.Split(targetPath, "/")
	pathHasWild := false
	for i, targetPathPart := range targetParts {
		if i > len(parts)-1 {
			return PartialMatch
		}
		treePathPart := parts[i]
		if pathPartIsWildcard(treePathPart) {
			pathHasWild = true
			continue
		}
		if treePathPart != targetPathPart {
			return NoMatch
		}
	}
	if pathHasWild {
		return WildMatch
	}
	return ExactMatch
}

// closest nodes in children in order of match level, greatest to lowest
func (tree *RouteTree) matchCandidates(targetPathPart string, targetMethod string) []*RouteTree {
	var out []*RouteTree
	for _, child := range tree.children {
		if child.method == targetMethod {
			if child.pathPart == targetPathPart {
				out = append([]*RouteTree{child}, out...)
			} else if child.isWild {
				out = append(out, child)
			}
		}
	}
	return out
}

// FindClosestMatchingNode searches the tree to find the node that best matches
// the requested path and method. It returns the closest node and a match level
// of ExactMatch, WildMatch, or NoMatch.
func (tree *RouteTree) FindClosestMatchingNode(targetPath string, targetMethod string) (*RouteTree, MatchLevel) {
	targetPath = strings.TrimRight(targetPath, "/")
	targetParts := strings.Split(targetPath, "/")
	var helper func(int, *RouteTree, bool) (int, *RouteTree, bool)
	helper = func(tpIndex int, current *RouteTree, isWildPath bool) (int, *RouteTree, bool) {
		if tpIndex >= len(targetParts) {
			return tpIndex - 1, current, isWildPath
		}
		tp := targetParts[tpIndex]
		bestCandidateDepth := tpIndex - 1
		bestCandidate := current
		candidates := current.matchCandidates(tp, targetMethod)
		for _, candidate := range candidates {
			i, node, _ := helper(tpIndex+1, candidate, isWildPath)
			if i > bestCandidateDepth {
				bestCandidateDepth = i
				bestCandidate = node
			}
			// set wild data
			if candidate.isWild {
				candidate.wildData = tp
			}
		}
		return bestCandidateDepth, bestCandidate, isWildPath || bestCandidate.isWild
	}
	i, closestNode, isWild := helper(0, tree, tree.isWild)
	if i == len(targetParts)-1 {
		if isWild {
			return closestNode, WildMatch
		}
		return closestNode, ExactMatch
	}
	return closestNode, NoMatch
}

// Deprecated:
// FindClosestMatchingNode breadth-first searches the tree for the closest matching node.
// It returns the closest node and the match level. Returns an error if tree is not the root.
func (tree *RouteTree) FindClosestMatchingNodeOLD(targetPath string, targetMethod string) (*RouteTree, MatchLevel, error) {
	if !tree.IsRoot() {
		return nil, NoMatch, errors.New("FindClosestMatchingNode called by non-root node")
	}
	if !strings.HasSuffix(targetPath, "/") {
		targetPath = targetPath + "/"
	}

	var matchAmtHelper = func(tree *RouteTree, path string, method string) MatchLevel {
		pathCmp := tree.compareWithPath(path)
		if method == tree.method {
			return pathCmp
		}
		return NoMatch
	}

	var closestNode *RouteTree = nil
	var closestMatchAmt MatchLevel = NoMatch
	queue := tree.children
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:] // dequeue current
		m := matchAmtHelper(current, targetPath, targetMethod)
		if m == ExactMatch {
			return current, m, nil
		} else {
			if m >= closestMatchAmt {
				closestMatchAmt = m
				closestNode = current
			}
			for _, child := range current.children {
				if matchAmtHelper(child, targetPath, targetMethod) >= closestMatchAmt {
					queue = append(queue, child)
				}
			}
		}
	}
	return closestNode, closestMatchAmt, nil
}

// FindClosestNotFoundHandler finds and returns the closest node with a non-nil
// notFoundHandler, searching up the tree from the current node.
func (tree *RouteTree) FindClosestNotFoundHandler() *RouteTree {
	if tree == nil {
		return nil
	}
	if tree.notFoundHandler != nil {
		return tree
	}
	return tree.parent.FindClosestNotFoundHandler()
}

func (tree *RouteTree) String() string {
	return tree.stringHelper(0)
}

func (tree *RouteTree) stringHelper(level int) string {
	if tree == nil {
		return ""
	}
	// node name
	spacer := ""
	for i := 0; i < level-1; i++ {
		spacer += "|   "
	}
	str := ""
	if level > 0 {
		str += spacer + "|\n"
		str += spacer + "└── "
		if tree.isWild {
			str += "{"
		}
		str += tree.pathPart
		if tree.isWild {
			str += "}"
		}
		str += "/"
	}
	if level == 1 {
		str += " [" + string(tree.method) + "]"
	}
	// more node info
	if tree.notFoundHandler != nil {
		str += " [404 Handler]"
	}
	// children
	for _, child := range tree.children {
		str += "\n" + child.stringHelper(level+1)
	}
	return str
}
