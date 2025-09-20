package test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"testing"
	"time"
)

type friendRequest struct {
	Id           int64  `json:"id"`
	UserId       string `json:"user_id,omitempty"`
	ReqUid       string `json:"req_uid,omitempty"`
	HandleResult int    `json:"handle_result,omitempty"`
}

type friendRequestList struct {
	List []friendRequest `json:"list"`
}

type friendEntry struct {
	FriendUid string `json:"friend_uid"`
}

type friendListData struct {
	List []friendEntry `json:"list"`
}

type mutualCountData struct {
	Count int64 `json:"count"`
}

type groupCreateResp struct {
	GroupId string `json:"group_id"`
}

type groupRequests struct {
	List []groupRequest `json:"list"`
}

type groupRequest struct {
	Id      int64  `json:"id"`
	GroupId string `json:"group_id"`
	UserId  string `json:"user_id"`
}

type groupMembers struct {
	List []groupMember `json:"List"`
}

type groupMember struct {
	UserId string `json:"user_id"`
}

type groupListData struct {
	List []groupInfo `json:"list"`
}

type groupInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type myCreatedGroupData struct {
	List []groupInfo `json:"list"`
}

const serverCommonErrorCode = 100001

type friendRequestContext struct {
	requester *userContext
	target    *userContext
	data      *friendRequest
}

func orderedUsers(users map[string]*userContext) []*userContext {
	labels := make([]string, 0, len(users))
	for label := range users {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	ordered := make([]*userContext, 0, len(labels))
	for _, label := range labels {
		ordered = append(ordered, users[label])
	}
	return ordered
}

func createPendingFriendRequest(t *testing.T, candidates []*userContext, exclude map[string]struct{}) *friendRequestContext {
	t.Helper()

	for i := range candidates {
		reqUser := candidates[i]
		if exclude != nil {
			if _, skip := exclude[reqUser.UserID]; skip {
				continue
			}
		}
		for j := range candidates {
			if i == j {
				continue
			}
			target := candidates[j]
			if exclude != nil {
				if _, skip := exclude[target.UserID]; skip {
					continue
				}
			}
			msg := fmt.Sprintf("auto friend %d", time.Now().UnixNano())
			pending, err := trySendFriendRequest(reqUser, target, msg)
			if err != nil {
				t.Fatalf("try friend request %s->%s: %v", reqUser.UserID, target.UserID, err)
			}
			if pending != nil {
				return &friendRequestContext{requester: reqUser, target: target, data: pending}
			}
		}
	}

	t.Fatalf("no available friend request slot; consider resetting test data")
	return nil
}

func createPendingFriendRequestWithParticipant(t *testing.T, participant *userContext, candidates []*userContext, exclude map[string]struct{}) *friendRequestContext {
	t.Helper()

	if exclude == nil {
		exclude = make(map[string]struct{})
	}
	for _, other := range candidates {
		if other.UserID == participant.UserID {
			continue
		}
		if _, skip := exclude[other.UserID]; skip {
			continue
		}
		msg := fmt.Sprintf("auto friend %d", time.Now().UnixNano())
		if pending, err := trySendFriendRequest(participant, other, msg); err != nil {
			t.Fatalf("try friend request %s->%s: %v", participant.UserID, other.UserID, err)
		} else if pending != nil {
			return &friendRequestContext{requester: participant, target: other, data: pending}
		}
		if pending, err := trySendFriendRequest(other, participant, msg); err != nil {
			t.Fatalf("try friend request %s->%s: %v", other.UserID, participant.UserID, err)
		} else if pending != nil {
			return &friendRequestContext{requester: other, target: participant, data: pending}
		}
	}

	t.Fatalf("participant %s has no available counterpart for new request", participant.UserID)
	return nil
}

func TestFriendWorkflow(t *testing.T) {
	users := ensureTestUsers(t)
	ordered := orderedUsers(users)

	first := createPendingFriendRequest(t, ordered, nil)

	outgoing, err := fetchFriendRequests(t, first.requester, 2)
	if err != nil {
		t.Fatalf("fetch outgoing friend requests: %v", err)
	}
	if findFriendRequest(outgoing, first.requester.UserID, first.target.UserID) == nil {
		t.Fatalf("pending request not visible in outgoing list for %s", first.requester.UserID)
	}

	incoming, err := fetchFriendRequests(t, first.target, 1)
	if err != nil {
		t.Fatalf("fetch incoming friend requests: %v", err)
	}
	if findFriendRequest(incoming, first.requester.UserID, first.target.UserID) == nil {
		t.Fatalf("pending request not visible in incoming list for %s", first.target.UserID)
	}

	approveFriendRequest(t, first.target, first.data.Id, 2)

	assertFriendshipExists(t, first.requester, first.target.UserID)
	assertFriendshipExists(t, first.target, first.requester.UserID)

	exclude := map[string]struct{}{first.requester.UserID: {}}
	second := createPendingFriendRequestWithParticipant(t, first.target, ordered, exclude)

	incomingSecond, err := fetchFriendRequests(t, second.target, 1)
	if err != nil {
		t.Fatalf("fetch incoming friend requests for second pair: %v", err)
	}
	if findFriendRequest(incomingSecond, second.requester.UserID, second.target.UserID) == nil {
		t.Fatalf("pending request not visible for second pair %s -> %s", second.requester.UserID, second.target.UserID)
	}

	approveFriendRequest(t, second.target, second.data.Id, 2)

	assertFriendshipExists(t, second.requester, second.target.UserID)
	assertFriendshipExists(t, second.target, second.requester.UserID)

	var third *userContext
	if second.requester.UserID == first.target.UserID {
		third = second.target
	} else {
		third = second.requester
	}

	count := mutualFriendCount(t, first.requester, third.UserID)
	if count < 1 {
		t.Fatalf("expected mutual count >=1 for %s and %s, got %d", first.requester.UserID, third.UserID, count)
	}
}

func TestGroupWorkflow(t *testing.T) {
	users := ensureTestUsers(t)
	owner := users["test04"]
	member := users["test05"]
	reviewer := users["test06"]

	groupID := createGroup(t, owner, fmt.Sprintf("integration-%d", time.Now().UnixNano()))

	sendGroupJoinRequest(t, member, groupID, "member join")

	requests, err := fetchGroupRequests(t, owner, groupID)
	if err != nil {
		if apiErr, ok := err.(*apiError); !ok || apiErr.Code != serverCommonErrorCode {
			t.Fatalf("fetch group requests: %v", err)
		}
		requests = nil
	}

	if len(requests) > 0 {
		req := requests[0]
		if req.GroupId != groupID {
			t.Fatalf("unexpected group id in first request: want %s got %s", groupID, req.GroupId)
		}
		approveGroupRequest(t, owner, req.Id, groupID, 2)
	}

	members := fetchGroupMembers(t, owner, groupID)
	if !containsUser(members, member.UserID) {
		t.Fatalf("member %s not found in group %s", member.UserID, groupID)
	}

	joined := fetchGroupList(t, member)
	if !containsGroup(joined, groupID) {
		t.Fatalf("group %s not visible for member %s", groupID, member.UserID)
	}

	created := fetchMyCreatedGroups(t, owner)
	if !containsGroup(created, groupID) {
		t.Fatalf("group %s not visible in owner created list", groupID)
	}

	sendGroupJoinRequest(t, reviewer, groupID, "reviewer join")
	reqs, err := fetchGroupRequests(t, owner, groupID)
	if err != nil {
		if apiErr, ok := err.(*apiError); !ok || apiErr.Code != serverCommonErrorCode {
			t.Fatalf("fetch group requests after reviewer join: %v", err)
		}
		reqs = nil
	}
	var reviewerReq *groupRequest
	for i := range reqs {
		if reqs[i].UserId == reviewer.UserID {
			reviewerReq = &reqs[i]
			break
		}
	}
	if reviewerReq == nil {
		t.Logf("group requests: %+v", reqs)
		t.Fatalf("expected reviewer join request")
	}
	approveGroupRequest(t, owner, reviewerReq.Id, groupID, 3)
}

func trySendFriendRequest(from, to *userContext, msg string) (*friendRequest, error) {
	payload := map[string]any{
		"user_uid": to.UserID,
		"req_msg":  msg,
		"req_time": time.Now().Unix(),
	}
	url := fmt.Sprintf("%s/v1/social/friend/putIn", socialAPIBase)
	if _, err := doJSONRequest(nil, "POST", url, from.Token, payload); err != nil {
		return nil, err
	}

	list, err := fetchFriendRequests(nil, from, 2)
	if err != nil {
		return nil, err
	}
	return findFriendRequest(list, from.UserID, to.UserID), nil
}

func fetchFriendRequests(tb testing.TB, user *userContext, direction int) ([]friendRequest, error) {
	if tb != nil {
		tb.Helper()
	}

	endpoint := fmt.Sprintf("%s/v1/social/friend/putIns", socialAPIBase)
	if direction > 0 {
		q := url.Values{}
		q.Set("direction", fmt.Sprintf("%d", direction))
		endpoint = endpoint + "?" + q.Encode()
	}
	envelope, err := doJSONRequest(tb, "GET", endpoint, user.Token, nil)
	if err != nil {
		return nil, err
	}

	var data friendRequestList
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		return nil, err
	}
	return data.List, nil
}

func approveFriendRequest(t *testing.T, user *userContext, requestID int64, result int32) {
	t.Helper()

	payload := map[string]any{
		"friend_req_id": requestID,
		"handle_result": result,
	}
	url := fmt.Sprintf("%s/v1/social/friend/putIn", socialAPIBase)
	if _, err := doJSONRequest(t, "PUT", url, user.Token, payload); err != nil {
		t.Fatalf("handle friend request: %v", err)
	}
}

func assertFriendshipExists(t *testing.T, user *userContext, friendID string) {
	t.Helper()

	url := fmt.Sprintf("%s/v1/social/friends", socialAPIBase)
	envelope, err := doJSONRequest(t, "GET", url, user.Token, nil)
	if err != nil {
		t.Fatalf("fetch friend list: %v", err)
	}

	var data friendListData
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		t.Fatalf("decode friend list: %v", err)
	}

	for _, item := range data.List {
		if item.FriendUid == friendID {
			return
		}
	}

	t.Fatalf("friend %s not present for user %s", friendID, user.UserID)
}

func mutualFriendCount(t *testing.T, user *userContext, otherID string) int64 {
	t.Helper()

	endpoint := fmt.Sprintf("%s/v1/social/friend/mutual/count?other_id=%s", socialAPIBase, url.QueryEscape(otherID))
	envelope, err := doJSONRequest(t, "GET", endpoint, user.Token, nil)
	if err != nil {
		t.Fatalf("mutual count: %v", err)
	}

	var data mutualCountData
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		t.Fatalf("decode mutual count: %v", err)
	}

	return data.Count
}

func findFriendRequest(list []friendRequest, requester, target string) *friendRequest {
	for i := range list {
		if list[i].ReqUid == requester && list[i].UserId == target {
			return &list[i]
		}
	}
	return nil
}

func createGroup(t *testing.T, owner *userContext, name string) string {
	t.Helper()

	payload := map[string]any{"name": name, "icon": "group.png"}
	url := fmt.Sprintf("%s/v1/social/group", socialAPIBase)
	envelope, err := doJSONRequest(t, "POST", url, owner.Token, payload)
	if err != nil {
		t.Fatalf("create group: %v", err)
	}

	var data groupCreateResp
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		t.Fatalf("decode group create: %v", err)
	}

	return data.GroupId
}

func sendGroupJoinRequest(t *testing.T, user *userContext, groupID, msg string) {
	t.Helper()

	payload := map[string]any{
		"group_id":    groupID,
		"req_msg":     msg,
		"req_time":    time.Now().Unix(),
		"join_source": 2,
	}
	url := fmt.Sprintf("%s/v1/social/group/putIn", socialAPIBase)
	if _, err := doJSONRequest(t, "POST", url, user.Token, payload); err != nil {
		t.Fatalf("group join request: %v", err)
	}
}

func fetchGroupRequests(tb testing.TB, owner *userContext, groupID string) ([]groupRequest, error) {
	if tb != nil {
		tb.Helper()
	}

	endpoint := fmt.Sprintf("%s/v1/social/group/putIns", socialAPIBase)
	if groupID != "" {
		endpoint = fmt.Sprintf("%s?group_id=%s", endpoint, url.QueryEscape(groupID))
	}
	envelope, err := doJSONRequest(tb, "GET", endpoint, owner.Token, nil)
	if err != nil {
		return nil, err
	}

	var data groupRequests
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		return nil, err
	}

	return data.List, nil
}

func approveGroupRequest(t *testing.T, owner *userContext, requestID int64, groupID string, result int32) {
	t.Helper()

	payload := map[string]any{
		"group_req_id":  requestID,
		"group_id":      groupID,
		"handle_result": result,
	}
	url := fmt.Sprintf("%s/v1/social/group/putIn", socialAPIBase)
	if _, err := doJSONRequest(t, "PUT", url, owner.Token, payload); err != nil {
		t.Fatalf("handle group request: %v", err)
	}
}

func fetchGroupMembers(t *testing.T, user *userContext, groupID string) []groupMember {
	t.Helper()

	endpoint := fmt.Sprintf("%s/v1/social/group/users?group_id=%s", socialAPIBase, url.QueryEscape(groupID))
	envelope, err := doJSONRequest(t, "GET", endpoint, user.Token, nil)
	if err != nil {
		t.Fatalf("fetch group members: %v", err)
	}

	var data groupMembers
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		t.Fatalf("decode group members: %v", err)
	}

	return data.List
}

func fetchGroupList(t *testing.T, user *userContext) []groupInfo {
	t.Helper()

	endpoint := fmt.Sprintf("%s/v1/social/groups", socialAPIBase)
	envelope, err := doJSONRequest(t, "GET", endpoint, user.Token, nil)
	if err != nil {
		t.Fatalf("fetch group list: %v", err)
	}

	var data groupListData
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		t.Fatalf("decode group list: %v", err)
	}

	return data.List
}

func fetchMyCreatedGroups(t *testing.T, user *userContext) []groupInfo {
	t.Helper()

	endpoint := fmt.Sprintf("%s/v1/social/groups/myCreated", socialAPIBase)
	envelope, err := doJSONRequest(t, "GET", endpoint, user.Token, nil)
	if err != nil {
		t.Fatalf("fetch my created groups: %v", err)
	}

	var data myCreatedGroupData
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		t.Fatalf("decode created groups: %v", err)
	}

	return data.List
}

func containsUser(list []groupMember, id string) bool {
	for _, item := range list {
		if item.UserId == id {
			return true
		}
	}
	return false
}

func containsGroup(list []groupInfo, groupID string) bool {
	for _, item := range list {
		if item.Id == groupID {
			return true
		}
	}
	return false
}
