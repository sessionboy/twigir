use crate::models::groups::member::MemberWithRole;

pub fn get_member(
  members: Option<Vec<MemberWithRole>>
) -> Option<MemberWithRole> {
  let mut member: Option<MemberWithRole> = None;
  if members.is_some() {
    let members_me = members.unwrap();
    member = Some(members_me.into_iter().next().unwrap());
  }
  member
}


