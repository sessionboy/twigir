
use json_gettext::JSONGetText;

pub fn langs_ctx() -> JSONGetText<'static> {
  static_json_gettext_build!(
    "zh_CN",
    // "en", "src/locales/langs/en.json",
    "zh_CN", "src/locales/langs/zh_CN.json"
  )
  .unwrap()
}