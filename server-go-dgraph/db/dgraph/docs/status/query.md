#### 1，查询贴文详情 (待定)
//  +","+ toString(u.auth_name)+","+ toString(u.username)+","+ toString(u.avatar_url)+","+ toString(u.verified) 
```js
MATCH (u:user)<-[]-(v:status)-[:to_status]->(t:status)
RETURN head(collect(distinct v.text)), u.name, head(collect(distinct t.text))
```
```js
LOOKUP ON status YIELD status.text AS text
| YIELD $-.VertexID AS vids, $-.text AS text
| GO FROM $-.vids OVER owner,to_status WHERE $-.vids == 1391787869177122816 YIELD $-.vids AS vids, $-.text AS text, $$.user.name as owner_name,$$.status.text as to_text 
| YIELD $-.vids AS id, $-.text, head(collect($-.owner_name)), collect($-.to_text)
```
```js
$b=go from 1391786075524960256 over owner yield
owner._src as id,
$^.status.text as text,
[owner._dst,$$.user.name,$$.user.username,$$.user.auth_name,$$.user.verified,$$.user.avatar_url] as owner,
$^.status.status_type as status_type,
$^.status.media_type as media_type,
$^.status.replies_count as replies_count,
$^.status.favorites_count as favorites_count,
$^.status.quotes_count as quotes_count,
$^.status.restatuses_count as restatuses_count,
$^.status.platform as platform,
$^.status.created_at as created_at;

[
[owner._src,$^.status.text,$^.status.status_type,$^.status.media_type,$^.status.replies_count,$^.status.favorites_count,
$^.status.quotes_count,$^.status.restatuses_count,$^.status.platform,$^.status.created_at] as to_status,
[owner._dst,$$.user.name,$$.user.username,$$.user.auth_name,$$.user.verified,$$.user.avatar_url] as owner
] as to_status
```
```js
go from 1391787869177122816 over to_user yield 
	NULL as status,
  NULL as owner, 
  [to_user._dst,$$.user.username] as to_user, 
  NULL as to_status
UNION DISTINCT 
go from 1391787869177122816 over owner yield
  [owner._src,$^.status.text,$^.status.status_type,$^.status.media_type,$^.status.replies_count,$^.status.favorites_count,
  $^.status.quotes_count,$^.status.restatuses_count,$^.status.platform,$^.status.created_at] as status,
  [owner._dst,$$.user.name,$$.user.username,$$.user.auth_name,$$.user.verified,$$.user.avatar_url] as owner,
  NULL as to_user, 
  NULL as to_status
UNION DISTINCT 
go from 1391787869177122816 over to_status yield to_status._dst as to_status_id |
go from $-.to_status_id over owner yield
  NULL as status, 
  NULL as owner,
  NULL as to_user, 
  [owner._dst] as to_status
```

#### 2，查询贴文回复列表
```js
GO FROM 1391786075524960256 OVER to_status REVERSELY 
where $$.status.status_type == 3 
YIELD 
	to_status._dst as id, 
  $$.status.text as text,
  $$.status.status_type as status_type,
  $$.status.media_type as media_type,
  $$.status.device as device,
  $$.status.platform as platform,
  $$.status.created_at as created_at,
  $$.status.favorites_count as favorites_count,
  $$.status.replies_count as replies_count|
ORDER BY favorites_count DESC, replies_count DESC |
LIMIT 0, 15
```