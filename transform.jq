.apps |
map(
  {
    url: (.packageName + "/"),
    contents: (
      (
        [
          .,
          (.localized // {})
          | to_entries 
          | map(select(.key | test("^(android$|en(-|_|$))")))
          | sort 
          | reverse 
          | map(.value)
        ] |
        flatten |
	reduce .[] as $item ({}; . * $item)
      ) |
      [.summary, .description] |
      unique_by(.) |
      join("\n\n")),
    title: (.localized["en-US"]?.name)
  }
) |
map(select(.title)) |
{
  input: {
    url_prefix: "https://f-droid.org/en/packages/",
    title_boost: "Ridiculous",
    files:.
  }
}
