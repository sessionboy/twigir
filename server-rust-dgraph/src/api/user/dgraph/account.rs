
// 查询name是否存在
pub const EXIST_NAME: &'static str = r#"query all($name: string) {
  name(func: eq(name,$name)) {
    uid 
  }
}"#;

// 查询username是否存在
pub const EXIST_USERNAME: &'static str = r#"query all($username: string) {
  username(func: eq(username,$username)) {
    uid 
  }
}"#;

// 查询phone_number是否存在
pub const EXIST_PHONE: &'static str = r#"query all($phone_number: string) {
  phone(func: eq(phone_number,$phone_number)) {
    uid 
  }
}"#;

// 查询email是否存在
pub const EXIST_EMAIL: &'static str = r#"query all($email: string) {
  email(func: eq(email,$email)) {
    uid 
  }
}"#;

// 查询用户密码
pub const QUERY_PASSWORD: &'static str = r#"query all($uid: string) {
  user(func: uid($uid)) {
    uid 
    password
  }
}"#;