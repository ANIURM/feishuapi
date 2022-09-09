package feishuapi

type SpaceType string
type Visibility string

const (
	Team   SpaceType = "team"
	Person SpaceType = "person"
)

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

type SpaceInfo struct {
	Name        string
	Description string
	SpaceId     string
	SpaceType   SpaceType
	Visibility  Visibility
}

// Create a new SpaceInfo
func NewSpaceInfo(data map[string]interface{}) *SpaceInfo {
	return &SpaceInfo{
		Name:        data["name"].(string),
		Description: data["description"].(string),
		SpaceId:     data["space_id"].(string),
		SpaceType:   data["space_type"].(SpaceType),
		Visibility:  data["visibility"].(Visibility),
	}
}

// Create a Knowledge Space
func (c AppClient) CreateKnowledgeSpace(name string, description string, user_access_token string) *SpaceInfo {
	body := make(map[string]string)
	body["name"] = name
	body["description"] = description

	headers := make(map[string]string)
	headers["Authorization"] = user_access_token

	info := c.Request("post", "open-apis/wiki/v2/spaces", nil, headers, body)

	return NewSpaceInfo(info)
}

type Node struct {
	NodeToken       string
	ParentNodeToken string
	Title           string
}

// Create a new Node
func NewNode(data map[string]interface{}) *Node {
	return &Node{
		NodeToken:       data["node_token"].(string),
		ParentNodeToken: data["parent_node_token"].(string),
		Title:           data["title"].(string),
	}
}

// Copy a node from SpaceId/NodeToken to TargetSpaceId/TargetParentToken
func (c AppClient) CopyNode(SpaceId string, NodeToken string, TargetSpaceId string, TargetParentToken string, Title ...string) *Node {
	body := make(map[string]string)
	body["target_parent_token"] = TargetParentToken
	body["target_space_id"] = TargetSpaceId
	if len(Title) != 0 {
		body["title"] = Title[0]
	}

	info := c.Request("post", "open-apis/wiki/v2/spaces/"+SpaceId+"/nodes/"+NodeToken+"/copy", nil, nil, body)

	return NewNode(info["node"].(map[string]interface{}))
}

type NodeInfo struct {
	NodeToken       string
	ObjToken        string
	ObjType         string
	ParentNodeToken string
	Title           string
}

// Create a new NodeInfo
func NewNodeInfo(data map[string]interface{}) *NodeInfo {
	return &NodeInfo{
		NodeToken:       data["node_token"].(string),
		ObjToken:        data["obj_token"].(string),
		ObjType:         data["obj_type"].(string),
		ParentNodeToken: data["parent_node_token"].(string),
		Title:           data["title"].(string),
	}
}

// Get All Nodes in target Space and under specific ParentNode(not necessary)
func (c AppClient) GetAllNodes(SpaceId string, ParentNodeToken ...string) []NodeInfo {
	var all_node []NodeInfo
	var l []interface{}

	if len(ParentNodeToken) != 0 {
		query := make(map[string]string)
		query["parent_node_token"] = ParentNodeToken[0]
		l = c.GetAllPages("get", "open-apis/wiki/v2/spaces/"+SpaceId+"/nodes", query, nil, nil, 10)
	} else {
		l = c.GetAllPages("get", "open-apis/wiki/v2/spaces/"+SpaceId+"/nodes", nil, nil, nil, 10)
	}

	for _, value := range l {
		all_node = append(all_node, *NewNodeInfo(value.(map[string]interface{})))
	}

	return all_node
}
