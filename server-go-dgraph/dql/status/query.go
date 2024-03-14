package status

import "fmt"

// 查询用户最近的发帖信息
var QueryStatusInfo = `query all($loggedUserid: string) {
  user(func: uid($loggedUserid)) {
    id:uid        
    last_publish_at
    status_count
  }
}`

// 查询用户最近引用发帖信息
type QuoteInfo struct {
	User   []QuoteLoggedUser `json:"user"`
	Status []QuoteStatus     `json:"status"`
}
type QuoteLoggedUser struct {
	Id            string `json:"id"`
	LastPublishAt string `json:"last_publish_at"`
	StatusCount   int    `json:"status_count"`
}
type QuoteStatus struct {
	Id         string `json:"id"`
	StatusType int    `json:"status_type"`
	CreatedAt  string `json:"created_at"`
	QuoteCount int    `json:"quote_count"`
	User       WithId `json:"user"`
}
type WithId struct {
	Id string `json:"id"`
}

var QueryQuoteInfo = `query all($loggedUserid: string, $statusid: string) {
  user(func: uid($loggedUserid)) {
    id:uid        
    last_publish_at
    status_count
  }
  status(func: uid($statusid)) {
    id:uid        
    user{
      id:uid
    }
    status_type
    quote_count
    created_at
  }
}`

// 查询用户最近的发帖信息
var UserStatusInfo = `query all($loggedUserid: string, $statusid: string) {
  user(func: uid($loggedUserid)) {
    id:uid        
    last_publish_at
    statuses_count
  }
  status(func: uid($statusid)) {
    id:uid    
    status_type    
    replies_count
  }
}`

// 查询转载的帖子
var QueryReStatus = `query all($statusid: string) {
  status(func: uid($statusid)) {
    id:uid   
    restatus_count
  }
}`

// 查询回复相关信息
type ReplyInfo struct {
	User   []LoggedUser `json:"user"`
	Status []Status     `json:"status"`
}
type LoggedUser struct {
	Id          string `json:"id"`
	LastReplyAt string `json:"last_reply_at"`
}
type Status struct {
	Id         string     `json:"id"`
	StatusType int        `json:"status_type"`
	CreatedAt  string     `json:"created_at"`
	ReplyCount int        `json:"reply_count"`
	User       LoggedUser `json:"user,omitempty"`
	ToStatus   ToStatus   `json:"to_status,omitempty"`
}
type User struct {
	Id string `json:"id"`
}
type ToStatus struct {
	Id string `json:"id"`
}

var GetReplyInfo = `query all($loggedUserid: string, $statusid: string) {
  user(func: uid($loggedUserid)) {
    id:uid        
    last_reply_at
  }
  status(func: uid($statusid)) {
    id:uid        
    user{
      id:uid
    }
    to_status{
      id:uid
    }
    reply_count
    created_at
  }
}`

// [删除贴文]查询用户和贴文信息
type DStatusInfo struct {
	User   []DLoggedUser `json:"user"`
	Status []DStatus     `json:"status"`
}
type DLoggedUser struct {
	Id          string `json:"id"`
	StatusCount int    `json:"status_count"`
}
type DStatus struct {
	Id         string `json:"id"`
	StatusType int    `json:"status_type,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	User       DUser  `json:"user,omitempty"`
	ToQuote    DQuote `json:"to_quote,omitempty"`
	ToReply    DReply `json:"to_reply,omitempty"`
}
type DUser struct {
	Id string `json:"id"`
}
type DQuote struct {
	Id         string `json:"id"`
	QuoteCount int    `json:"quote_count"`
}
type DReply struct {
	Id         string `json:"id"`
	ReplyCount int    `json:"reply_count"`
}

var QueryDeleteStatusInfo = `query all($loggedUserid: string, $statusid: string) {
  user(func: uid($loggedUserid)) {
    id:uid        
    status_count
  }
  status(func: uid($statusid)) {
    id:uid        
    status_type
    created_at
    user{
      id:uid
    }
    to_quote{
      id:uid
      quote_count
    }
    to_reply{
      id:uid
      reply_count
    }
  }
}`

// 查询贴文详情
var QueryStatus = `query all($statusid: string,$loggedUserid: string) {
  status(func: uid($statusid)) {
    id:uid        
    text
    status_type
    media_type
    user{
      name
      username
      verify_name
      avatar_url
      verified
    }
    to_user{
      name
      username
      verify_name
    }
    urls
    images{
      url
      source
      platform
    }
    hashtags{
      id:uid
      name
      description
    }
    video{
      url
      source
      platform
    }
    reply_count
    quote_count
    restatus_count
    favorite_count
    platform

    favorite_cnt as fcnt:count(~favorites @filter(uid($loggedUserid)))
    favorited: math(favorite_cnt == 1)

    restatus_cnt as rcnt:count(~restatuses @filter(uid($loggedUserid)))
    restatused: math(restatus_cnt == 1)

    created_at
    to_quote{
      id:uid        
      text
      status_type
      media_type
      user{
        name
        username
        verify_name
        avatar_url
        verified
      }
      to_user{
        name
        username
        verify_name
      }
      urls
      images{
        url
        source
        platform
      }
      hashtags{
        id:uid
        name
        description
      }
      video{
        url
        source
        platform
      }
      created_at
    }
  }
}`

// 查询贴文回复列表
var QueryStatusReplies = `query all(
  $statusid: string,
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  status(func: uid($statusid)) {
    id:uid
    edges: ~to_reply (
      orderdesc: verified, 
      orderdesc: reply_count, 
      orderdesc: created_at, 
      first: $first, 
      after:$after
    ) {
      id:uid        
      text
      media_type
      user{
        name
        username
        verify_name
        avatar_url
        verified
      }
      to_user{
        name
        username
        verify_name
      }
      urls
      images{
        url
        source
        platform
      }
      hashtags{
        id:uid
        name
        description
      }
      video{
        url
        source
        platform
      }
      reply_count
      quote_count
      restatus_count
      favorite_count
      platform

      favorite_cnt as fcnt:count(~favorites @filter(uid($loggedUserid)))
      favorited: math(favorite_cnt == 1)

      restatus_cnt as rcnt:count(~restatuses @filter(uid($loggedUserid)))
      restatused: math(restatus_cnt == 1)

      created_at
      to_quote{
        id:uid        
        text
        status_type
        media_type
        user{
          name
          username
          verify_name
          avatar_url
          verified
        }
        to_user{
          name
          username
          verify_name
        }
        urls
        images{
          url
          source
          platform
        }
        hashtags{
          id:uid
          name
          description
        }
        video{
          url
          source
          platform
        }
        created_at
      }
    }
  }
}`

// 为你推荐
var QueryRecommendStatuses = `query all(
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  statuses(
    func: type(Status), 
    orderdesc: verify_int, 
    orderdesc: created_at, 
    orderdesc: reply_count,
    first: $first, 
    after:$after
  ) @filter(lt(status_type,2)) {
    id:uid        
    text
    media_type
    user{
      name
      username
      verify_name
      avatar_url
      verified
    }
    to_user{
      name
      username
      verify_name
    }
    urls
    images{
      url
      source
      platform
    }
    hashtags{
      id:uid
      name
      description
    }
    video{
      url
      source
      platform
    }
    reply_count
    quote_count
    restatus_count
    favorite_count
    platform

    favorite_cnt as fcnt:count(~favorites @filter(uid($loggedUserid)))
    favorited: math(favorite_cnt == 1)

    restatus_cnt as rcnt:count(~restatuses @filter(uid($loggedUserid)))
    restatused: math(restatus_cnt == 1)

    created_at
    to_quote{
      id:uid        
      text
      status_type
      media_type
      user{
        name
        username
        verify_name
        avatar_url
        verified
      }
      to_user{
        name
        username
        verify_name
      }
      urls
      images{
        url
        source
        platform
      }
      hashtags{
        id:uid
        name
        description
      }
      video{
        url
        source
        platform
      }
      created_at
    }
  }
}`

// 为你推荐媒体贴文
var QueryRecommendMediaStatuses = `query all(
  $loggedUserid: string,
  $media_type: string,
  $first: string, 
  $after: string
) {
  statuses(
    func: type(Status), 
    orderdesc: verify_int, 
    orderdesc: created_at, 
    orderdesc: favorite_count,
    first: $first, 
    after:$after
  ) @filter(lt(status_type,2) AND eq(media_type,$media_type)) {
    id:uid        
    text
    media_type
    user{
      name
      username
      verify_name
      avatar_url
      verified
    }
    to_user{
      name
      username
      verify_name
    }
    urls
    images{
      url
      source
      platform
    }
    hashtags{
      id:uid
      name
      description
    }
    video{
      url
      source
      platform
    }
    reply_count
    quote_count
    restatus_count
    favorite_count
    platform

    favorite_cnt as fcnt:count(~favorites @filter(uid($loggedUserid)))
    favorited: math(favorite_cnt == 1)

    restatus_cnt as rcnt:count(~restatuses @filter(uid($loggedUserid)))
    restatused: math(restatus_cnt == 1)

    created_at
    to_quote{
      id:uid        
      text
      status_type
      media_type
      user{
        name
        username
        verify_name
        avatar_url
        verified
      }
      to_user{
        name
        username
        verify_name
      }
      urls
      images{
        url
        source
        platform
      }
      hashtags{
        id:uid
        name
        description
      }
      video{
        url
        source
        platform
      }
      created_at
    }
  }
}`

// 全局搜索
func SearchQuery(search_type string) string {
	users := ""
	sort := `
    orderdesc: verify_int, 
    orderdesc: created_at, 
    orderdesc: favorite_count
  `
	filter := "@filter(lt(status_type,2))"

	if search_type == "latest" {
		// 按最新排序
		sort = `orderdesc: created_at`
	} else if search_type == "image" {
		// 图片
		filter = "@filter(lt(status_type,2) AND eq(media_type,1))"
	} else if search_type == "video" {
		// 视频
		filter = "@filter(lt(status_type,2) AND eq(media_type,2))"
	} else if search_type == "user" {
		return `
      query all(
        $keyword: string,
        $loggedUserid: string,
        $first: string, 
        $after: string
      ){
        users(
          func: type(User),  
          orderdesc: verify_int, 
          orderdesc: created_at, 
          orderdesc: followers_count,
          first: $first,
          after: $after
        ) @filter(anyofterms(name, $keyword) OR anyofterms(username, $keyword) OR anyofterms(bio, $keyword)) {
          name
          username
          verify_name
          verified
          avatar_url          
          bio
          cnt as cnt:count(~follows @filter(uid($loggedUserid)))
          following: math(cnt == 1)   
        }
      }
    `
	} else {
		// all 综合查询
		users = `
      users(
        func: type(User),  
        orderdesc: verify_int, 
        orderdesc: created_at, 
        orderdesc: followers_count,
        first: 3
      ) @filter(anyofterms(name, $keyword) OR anyofterms(username, $keyword) OR anyofterms(bio, $keyword)) {
        name
        username
        verify_name
        avatar_url
        verified
        bio
        cnt as cnt:count(~follows @filter(uid($loggedUserid)))
        following: math(cnt == 1)   
      }
    `
	}
	q := fmt.Sprintf(`query all(
    $keyword: string,
    $loggedUserid: string,
    $first: string, 
    $after: string
  ) {
    %s
    statuses(
      func: alloftext(text, $keyword), 
      first: $first, 
      after:$after,
      %s  
    ) %s {
      id:uid        
      text
      media_type
      user{
        name
        username
        verify_name
        avatar_url
        verified
      }
      to_user{
        name
        username
        verify_name
      }
      urls
      images{
        url
        source
        platform
      }
      hashtags{
        id:uid
        name
        description
      }
      video{
        url
        source
        platform
      }
      reply_count
      quote_count
      restatus_count
      favorite_count
      platform
  
      favorite_cnt as fcnt:count(~favorites @filter(uid($loggedUserid)))
      favorited: math(favorite_cnt == 1)
  
      restatus_cnt as rcnt:count(~restatuses @filter(uid($loggedUserid)))
      restatused: math(restatus_cnt == 1)
  
      created_at
      to_quote{
        id:uid        
        text
        status_type
        media_type
        user{
          name
          username
          verify_name
          avatar_url
          verified
        }
        to_user{
          name
          username
          verify_name
        }
        urls
        images{
          url
          source
          platform
        }
        hashtags{
          id:uid
          name
          description
        }
        video{
          url
          source
          platform
        }
        created_at
      }
    }
  }`, users, sort, filter)

	return q
}
