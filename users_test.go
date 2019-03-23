package qiita

//func TestClient_GetFollowees(t *testing.T) {
//	tests := []struct {
//		desc         string
//		userID       string
//		page         int
//		perPage      int
//		responseFile string
//
//		expectedMethod      string
//		expectedRequestPath string
//		expectedRawQuery    string
//		expectedErrString   string
//		expectedPage        int
//		expectedPerPage     int
//		expectedFirstPage   int
//		expectedLastPage    int
//		expectedTotalCount  int
//		expectedUsersLen    int
//	}{
//		{
//			desc:         "success",
//			userID:       "muiscript",
//			page:         2,
//			perPage:      2,
//			responseFile: "users_muiscript_followees?page=2&per_page=2",
//
//			expectedMethod:      http.MethodGet,
//			expectedRequestPath: "/users/muiscript/followees",
//			expectedRawQuery:    "page=2&per_page=2",
//			expectedPage:        2,
//			expectedPerPage:     2,
//			expectedFirstPage:   1,
//			expectedLastPage:    3,
//			expectedTotalCount:  5,
//			expectedUsersLen:    2,
//		},
//		{
//			desc:         "success_page_larger_than_last",
//			userID:       "muiscript",
//			page:         10,
//			perPage:      2,
//			responseFile: "users_muiscript_followees?page=10&per_page=2",
//
//			expectedMethod:      http.MethodGet,
//			expectedRequestPath: "/users/muiscript/followees",
//			expectedRawQuery:    "page=10&per_page=2",
//			expectedPage:        10,
//			expectedPerPage:     2,
//			expectedFirstPage:   1,
//			expectedLastPage:    3,
//			expectedTotalCount:  5,
//			expectedUsersLen:    0,
//		},
//		{
//			desc:         "failure_page_less_than_100",
//			userID:       "muiscript",
//			page:         0,
//			perPage:      2,
//			responseFile: "users_muiscript_followees?page=0&per_page=2",
//
//			expectedMethod:      http.MethodGet,
//			expectedRequestPath: "/users/muiscript/followees",
//			expectedRawQuery:    "page=0&per_page=2",
//			expectedErrString:   "page parameter should be",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.desc, func(t *testing.T) {
//			server := newTestServer(t, tt.responseFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
//			defer server.Close()
//
//			serverURL, err := url.Parse(server.URL)
//			assert.Nil(t, err)
//			cli := &Client{
//				URL:        serverURL,
//				HTTPClient: server.Client(),
//				Logger:     log.New(ioutil.Discard, "", 0),
//			}
//
//			usersResp, err := cli.GetUserFollowees(context.Background(), tt.userID, tt.page, tt.perPage)
//			if tt.expectedErrString == "" {
//				if !assert.Nil(t, err) {
//					t.FailNow()
//				}
//
//				assert.Equal(t, tt.expectedPage, usersResp.Page)
//				assert.Equal(t, tt.expectedPerPage, usersResp.PerPage)
//				assert.Equal(t, tt.expectedFirstPage, usersResp.FirstPage)
//				assert.Equal(t, tt.expectedLastPage, usersResp.LastPage)
//				assert.Equal(t, tt.expectedTotalCount, usersResp.TotalCount)
//				assert.Equal(t, tt.expectedUsersLen, len(usersResp.Users))
//			} else {
//				if !assert.NotNil(t, err) {
//					t.FailNow()
//				}
//
//				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString))
//			}
//		})
//	}
//}
//
//func TestClient_GetFollowers(t *testing.T) {
//	tests := []struct {
//		desc         string
//		userID       string
//		page         int
//		perPage      int
//		responseFile string
//
//		expectedMethod      string
//		expectedRequestPath string
//		expectedRawQuery    string
//		expectedErrString   string
//		expectedPage        int
//		expectedPerPage     int
//		expectedFirstPage   int
//		expectedLastPage    int
//		expectedTotalCount  int
//		expectedUsersLen    int
//	}{
//		{
//			desc:         "success",
//			userID:       "muiscript",
//			page:         2,
//			perPage:      2,
//			responseFile: "users_muiscript_followers?page=2&per_page=2",
//
//			expectedMethod:      http.MethodGet,
//			expectedRequestPath: "/users/muiscript/followers",
//			expectedRawQuery:    "page=2&per_page=2",
//			expectedPage:        2,
//			expectedPerPage:     2,
//			expectedFirstPage:   1,
//			expectedLastPage:    6,
//			expectedTotalCount:  11,
//			expectedUsersLen:    2,
//		},
//		{
//			desc:         "success_page_larger_than_last",
//			userID:       "muiscript",
//			page:         10,
//			perPage:      2,
//			responseFile: "users_muiscript_followers?page=10&per_page=2",
//
//			expectedMethod:      http.MethodGet,
//			expectedRequestPath: "/users/muiscript/followers",
//			expectedRawQuery:    "page=10&per_page=2",
//			expectedPage:        10,
//			expectedPerPage:     2,
//			expectedFirstPage:   1,
//			expectedLastPage:    6,
//			expectedTotalCount:  11,
//			expectedUsersLen:    0,
//		},
//		{
//			desc:         "failure_page_less_than_100",
//			userID:       "muiscript",
//			page:         0,
//			perPage:      2,
//			responseFile: "users_muiscript_followers?page=0&per_page=2",
//
//			expectedMethod:      http.MethodGet,
//			expectedRequestPath: "/users/muiscript/followers",
//			expectedRawQuery:    "page=0&per_page=2",
//			expectedErrString:   "page parameter should be",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.desc, func(t *testing.T) {
//			server := newTestServer(t, tt.responseFile, tt.expectedMethod, tt.expectedRequestPath, tt.expectedRawQuery)
//			defer server.Close()
//
//			serverURL, err := url.Parse(server.URL)
//			assert.Nil(t, err)
//			cli := &Client{
//				URL:        serverURL,
//				HTTPClient: server.Client(),
//				Logger:     log.New(ioutil.Discard, "", 0),
//			}
//
//			usersResp, err := cli.GetUserFollowers(context.Background(), tt.userID, tt.page, tt.perPage)
//			if tt.expectedErrString == "" {
//				if !assert.Nil(t, err) {
//					t.FailNow()
//				}
//
//				assert.Equal(t, tt.expectedPage, usersResp.Page)
//				assert.Equal(t, tt.expectedPerPage, usersResp.PerPage)
//				assert.Equal(t, tt.expectedFirstPage, usersResp.FirstPage)
//				assert.Equal(t, tt.expectedLastPage, usersResp.LastPage)
//				assert.Equal(t, tt.expectedTotalCount, usersResp.TotalCount)
//				assert.Equal(t, tt.expectedUsersLen, len(usersResp.Users))
//			} else {
//				if !assert.NotNil(t, err) {
//					t.FailNow()
//				}
//
//				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString))
//			}
//		})
//	}
//}
//
//func TestClient_GetUser(t *testing.T) {
//	tests := []struct {
//		desc         string
//		id           string
//		responseFile string
//
//		expectedMethod         string
//		expectedRequestPath    string
//		expectedErrString      string
//		expectedID             string
//		expectedPermanentID    int
//		expectedGithubID       string
//		expectedPostsCount     int
//		expectedFollowersCount int
//	}{
//		{
//			desc:         "success",
//			id:           "muiscript",
//			responseFile: "users_muiscript",
//
//			expectedMethod:         http.MethodGet,
//			expectedRequestPath:    "/users/muiscript",
//			expectedID:             "muiscript",
//			expectedPermanentID:    159260,
//			expectedGithubID:       "muiscript",
//			expectedPostsCount:     14,
//			expectedFollowersCount: 11,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.desc, func(t *testing.T) {
//			server := newTestServer(t, tt.responseFile, tt.expectedMethod, tt.expectedRequestPath, "")
//			defer server.Close()
//
//			serverURL, err := url.Parse(server.URL)
//			assert.Nil(t, err)
//			cli := &Client{
//				URL:        serverURL,
//				HTTPClient: server.Client(),
//				Logger:     log.New(ioutil.Discard, "", 0),
//			}
//
//			user, err := cli.GetUser(context.Background(), tt.id)
//			if tt.expectedErrString == "" {
//				if !assert.Nil(t, err) {
//					t.FailNow()
//				}
//
//				assert.Equal(t, tt.expectedID, user.ID)
//				assert.Equal(t, tt.expectedPermanentID, user.PermanentID)
//				assert.Equal(t, tt.expectedGithubID, user.GithubID)
//				assert.Equal(t, tt.expectedGithubID, user.GithubID)
//				assert.Equal(t, tt.expectedPostsCount, user.PostsCount)
//				assert.Equal(t, tt.expectedFollowersCount, user.FollowersCount)
//			} else {
//				if !assert.NotNil(t, err) {
//					t.FailNow()
//				}
//
//				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString))
//			}
//
//		})
//	}
//}
//
//func TestClient_IsFollowingUser(t *testing.T) {
//	tests := []struct {
//		desc         string
//		targetUserID string
//		responseFile string
//
//		expectedMethod      string
//		expectedRequestPath string
//		expectedIsFollowing bool
//		expectedErrString   string
//	}{
//		{
//			desc:         "success_following",
//			targetUserID: "mizchi",
//			responseFile: "users_mizchi_following",
//
//			expectedMethod:      http.MethodGet,
//			expectedRequestPath: "/users/mizchi/following",
//			expectedIsFollowing: true,
//		},
//		{
//			desc:         "success_not_following",
//			targetUserID: "yaotti",
//			responseFile: "users_yaotti_following",
//
//			expectedMethod:      http.MethodGet,
//			expectedRequestPath: "/users/yaotti/following",
//			expectedIsFollowing: false,
//		},
//		{
//			desc:         "failure_no_token",
//			targetUserID: "mizchi",
//			responseFile: "users_mizchi_following-no_token",
//
//			expectedMethod:      http.MethodGet,
//			expectedRequestPath: "/users/mizchi/following",
//			expectedErrString:   "unauthorized",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.desc, func(t *testing.T) {
//			server := newTestServer(t, tt.responseFile, tt.expectedMethod, tt.expectedRequestPath, "")
//			defer server.Close()
//
//			serverURL, err := url.Parse(server.URL)
//			if !assert.Nil(t, err) {
//				t.FailNow()
//			}
//			cli := &Client{
//				URL:        serverURL,
//				HTTPClient: server.Client(),
//				Logger:     log.New(ioutil.Discard, "", 0),
//			}
//
//			isFollowing, err := cli.IsFollowingUser(context.Background(), tt.targetUserID)
//			if tt.expectedErrString == "" {
//				if !assert.Nil(t, err) {
//					t.FailNow()
//				}
//
//				assert.Equal(t, tt.expectedIsFollowing, isFollowing)
//			} else {
//				if !assert.NotNil(t, err) {
//					t.FailNow()
//				}
//
//				assert.True(t, strings.Contains(err.Error(), tt.expectedErrString))
//			}
//
//		})
//	}
//}
