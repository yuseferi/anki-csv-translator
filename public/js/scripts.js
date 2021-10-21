(function ($) {
    $(document).ready(function () {
        adjustOptions();
        $("#longman,#linguee").change(function () {
            adjustOptions();
        });

    });

    function adjustOptions() {
        if ($("#linguee").is(':checked')) {
            $("#iframe").attr("disabled", false).prop('checked', true).click();
            $("#full_html").attr("disabled", true);
        }
        if ($("#longman").is(':checked')) {
            $("#full_html").attr("disabled", false).prop('checked', true).click();
            $("#iframe").attr("disabled", false);
        }
    }
})(jQuery);