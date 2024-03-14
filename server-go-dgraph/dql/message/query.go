package message

// 查询是否存在对话
var HasConversation = `query all(
	$cid0: string, 
	$cid1: string
) {
  conversation(func: type(Conversation)) @filter(eq(conversation_id,$cid0) OR eq(conversation_id,$cid1)) {
    id:uid  
		conversation_id
  }
}`

// 查询对话
type Conversation struct {
	Conversation []ConversationData `json:"conversation"`
}
type ConversationData struct {
	Id             string `json:"id"`
	ConversationId string `json:"conversation_id"`
}

var ConversationWithId = `query all($conversationid: string) {
	conversation(func: uid($conversationid)) {
		id:uid
		conversation_id
	}
}`

// 检查用户
type ConversationUser struct {
	User []ConversationUserData `json:"user"`
}
type ConversationUserData struct {
	Id             string `json:"id"`
	ChatUnFollow   bool   `json:"chat_unFollow"`
	ChatUnFollowMe bool   `json:"chat_unFollowMe"`
	ChatUnVerified bool   `json:"chat_unVerified"`
	ChatBlacklist  bool   `json:"chat_blacklist"`
	Following      bool   `json:"following"`
	Followme       bool   `json:"followme"`
	Blacklist      bool   `json:"blacklist"`
}

var ConversationAuth = `query all(
	$userid: string,
	$loggedUserid: string
) {
  user(func: uid($userid)) {
    id:uid  
		chat_unFollow
    chat_unFollowMe
    chat_unVerified
    chat_blacklist

		fi as fi:count(follows @filter(uid($loggedUserid)))
    following: math(fi == 1)

		fm as fm:count(~follows @filter(uid($loggedUserid)))
    followme: math(fm == 1)

		bl as bl:count(blacklists @filter(uid($loggedUserid)))
    blacklist: math(bl == 1)
  }
}`

// 查询用户私信对话列表
func GetConversationsDql(
	loggedUserid string,
	first string,
	after string,
	keyword string,
) string {
	if len(keyword) != 0 {
		return `query all(
			$keyword: string, 
			$loggedUserid: string,
			$first: string, 
			$after: string
		) {
			users(func:allofterms(name,$keyword)) @filter(NOT uid($loggedUserid)) {
				id:uid
				conversations: ~users @filter(type(Conversation) AND uid_in(users,$loggedUserid)) (first: $first, after:$after, orderdesc: last_publish_at) {
					id:uid
					conversation_id
					users @filter(NOT uid($loggedUserid)) {
						id:uid  
						name
						username
						verify_name
						avatar_url
						verified
					}
					messages (first: 1, orderdesc: created_at) {
						msg
						sender{
							name
						}
					}
				}
			}
		}`
	}
	return `query all(
		$loggedUserid: string,
		$first: string, 
		$after: string
	) {
		conversations(
			func: type(Conversation),
			first: $first, 
			after:$after,
			orderdesc: last_publish_at
		) @filter(uid_in(users,$loggedUserid)) {
			id:uid
			conversation_id
			users @filter(NOT uid($loggedUserid)) {
				id:uid  
				name
				username
				verify_name
				avatar_url
				verified
			}
			messages (first: 1, orderdesc: created_at) {
				msg
				sender{
					name
				}
			}
		}
	}`
}

var QueryConversations = `query all(
	$keyword: string, 
	$loggedUserid: string,
  $first: string, 
  $after: string
) {
	users(func: allofterms(name,$keyword)){
		U as uid
	}
	conversations(
		func: type(Conversation),
		first: $first, 
		after:$after
	) @filter(uid_in(users,$loggedUserid) AND anyofterms(conversation_id,U) ) {
		id:uid
		conversation_id
		users @filter(NOT uid($loggedUserid)) {
			id:uid  
			name
			username
			verify_name
			avatar_url
			verified
		}
		messages (first: 1, orderdesc: created_at) {
			msg
			sender{
				name
			}
		}
	}
}`

// 查询消息
var QueryMessages = `query all(
	$conversationid: string,
	$first: string, 
  $after: string
) {
	conversations(func: uid($conversationid)) {
		id:uid
		messages (first: $first, after: $after, orderdesc: created_at) {
			id:uid  
			sender{
				id:uid  
				name
				username
				verify_name
				avatar_url
				verified
			}
			recipient{
				id:uid  
			}
			msg
			msg_type
			media_url
			created_at      
		}
	}
}`

// 查询单条消息
var QueryMessage = `query all($messageid: string) {
	message(func: uid($messageid)) {
		id:uid  
		sender{
			id:uid  
			name
			username
			verify_name
			avatar_url
			verified
		}
		recipient{
			id:uid  
		}
		msg
		msg_type
		media_url
		created_at 
	}
}`
