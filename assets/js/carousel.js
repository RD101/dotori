$('#recentlyCarousel').on('slide.bs.carousel', function (e) {
    var $e = $(e.relatedTarget);    // object
    var idx = $e.index();           // object Index
    console.log(idx)
    var itemsPerSlide = 4;
    var totalItems = $('div[name=recentlyCarousel-item]').length;
    
    if (idx >= totalItems-(itemsPerSlide-1)) {
        var it = itemsPerSlide - (totalItems - idx);
        console.log(it)
        for (var i=0; i<it; i++) {
            // append slides to end
            if (e.direction=="left") {
                console.log("left")
                $('div[name=recentlyCarousel-item]').eq(i).appendTo($('div[name=recentlyCarousel-inner]'));
            }
            else {
                console.log("none")
                $('div[name=recentlyCarousel-item]').eq(0).appendTo($('div[name=recentlyCarousel-inner]'));
            }
        }
    }
});

$('#recentlyCarousel').carousel({interval: 2000});

$('#topUsingCarousel').on('slide.bs.carousel', function (e) {
    var $e = $(e.relatedTarget);
    var idx = $e.index();
    var itemsPerSlide = 4;
    var totalItems = $('div[name=topUsingCarousel-item]').length;
    
    if (idx >= totalItems-(itemsPerSlide-1)) {
        var it = itemsPerSlide - (totalItems - idx);
        for (var i=0; i<it; i++) {
            // append slides to end
            if (e.direction=="left") {
                $('div[name=topUsingCarousel-item]').eq(i).appendTo($('div[name=topUsingCarousel-inner]'));
            }
            else {
                $('div[name=topUsingCarousel-item]').eq(0).appendTo($('div[name=topUsingCarousel-inner]'));
            }
        }
    }
});

$('#topUsingCarousel').carousel({interval: 2000});