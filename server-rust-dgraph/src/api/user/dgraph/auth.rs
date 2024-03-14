
// 注册 - 查询名称、主名、手机号是否已注册
pub const USER_EXIST: &'static str = r#"query all(
  $name: string,
  $username: string, 
  $phone: string
) {
  name(func: eq(name,$name)) {
    uid        
  }
  username(func: eq(username,$username)) {
    uid
  }
  phone(func: eq(phone_number,$phone)) {
    uid
  }
}"#;

// 账号登录
pub const LOGIN_BY_USERNAME: &'static str = r#"query all($username: string) {
  data(func: eq(username, $username)) {
    uid
    name
    username
    password
    role
    avatar_url
    description
    lang
    is_verified
    verified{
      uid
      name
      description
    }
  }
}"#;

// 手机号登录
pub const LOGIN_BY_PHONE: &'static str = r#"query all($phone_number: string) {
  data(func: eq(phone_number, $phone_number)) {
    uid
    name
    username
    password
    role
    avatar_url
    description
    lang
    is_verified
    verified{
      uid
      name
      description
    }
  }
}"#;

// 邮箱登录
pub const LOGIN_BY_EMAIL: &'static str = r#"query all($email: string) {
  data(func: eq(email, $email)) {
    uid
    name
    username
    password
    role
    avatar_url
    description
    lang
    is_verified
    verified{
      name
      description
    }
  }
}"#;