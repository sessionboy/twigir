
// 查询小组名称是否已存在
pub const GROUP_NAME_EXIST: &'static str = r#"query all($group_name: string) {
  data(func: eq(group_name,$group_name)) {
    uid
  }
}"#;

