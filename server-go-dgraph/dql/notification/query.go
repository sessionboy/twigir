package status

// 查询回复通知列表
var Notification_Replies = `query all(
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  notifications(
		func: type(Notification),
		orderdesc: created_at,
		first: $first,
		after:$after
	) @filter(uid_in(recipient,$loggedUserid) AND eq(action,0)) {
		id:uid        
		sender{
			id:uid        
			name
			username
			verify_name
			avatar_url
			verified
		}
		created_at
		target_type
		target{
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
			to_reply{
				id:uid        
				text
				user{
					name
					verify_name
				}
			}
		}
  }
}`

// 查询点赞通知列表
var Notification_Favorites = `query all(
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  notifications(
		func: type(Notification),
		orderdesc: created_at, 
		first: $first, 
    after:$after
	) @filter(uid_in(recipient,$loggedUserid) AND eq(action,1)) {
		id:uid        
		sender{
			id:uid        
			name
			username
			verify_name
			avatar_url
			verified
		}
		created_at
		target_type
		target{
			id:uid        
			text
			to_user{
				name
				verify_name
			}
			hashtags{
				id:uid
				name
				description
			}
		}
  }
}`

// 查询转帖通知列表
var Notification_Restatus = `query all(
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  notifications(
		func: type(Notification),
		orderdesc: created_at, 
		first: $first, 
    after:$after
	) @filter(uid_in(recipient,$loggedUserid) AND eq(action,2)) {
		id:uid        
		sender{
			id:uid        
			name
			username
			verify_name
			avatar_url
			verified
		}
		created_at
		target_type
		target{
			id:uid        
			text
			to_user{
				name
				verify_name
			}
			hashtags{
				id:uid
				name
				description
			}
		}
  }
}`

// 查询引用通知列表
var Notification_Quotes = `query all(
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  notifications(
		func: type(Notification),
		orderdesc: created_at, 
		first: $first, 
    after:$after
	) @filter(uid_in(recipient,$loggedUserid) AND eq(action,3)) {
		id:uid        
		created_at
		target_type
		target{
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
				user{
					name
					verify_name
				}
			}
		}
  }
}`

// 查询提及通知列表
var Notification_Mentions = `query all(
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  notifications(
		func: type(Notification),
		orderdesc: created_at, 
		first: $first, 
    after:$after
	) @filter(uid_in(recipient,$loggedUserid) AND eq(action,4)) {
		id:uid        
		created_at
		target_type
		target{
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
		}
  }
}`

// 查询用户通知列表
var Notification_Users = `query all(
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  notifications(
		func: type(Notification),
		orderdesc: created_at, 
		first: $first, 
    after:$after
	) @filter(uid_in(recipient,$loggedUserid) AND eq(action,5)) {
		id:uid        
		action
		created_at
		target_type
		target{
			id:uid        
			name
			username
			verify_name
			avatar_url
			verified
		}
  }
}`

// 查询广播通知列表
var Notification_Broadcast = `query all(
  $loggedUserid: string,
  $first: string, 
  $after: string
) {
  notifications(
		func: type(Notification),
		orderdesc: created_at, 
		first: $first, 
    after:$after
	) @filter(uid_in(recipient,$loggedUserid) AND eq(action,6)) {
		id:uid        
		action
		msg
		created_at
  }
}`
