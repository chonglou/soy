$(function() {
  $("div.markdown").each(function(e) {
    $(this).html(marked($(this).text()));
  });
});
