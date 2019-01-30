
    $(".navigation").onePageNav({
        currentClass: 'current',
        changeHash: false,
        scrollSpeed: 700,
        scrollThreshold: 0.5,
        easing: 'easeInOutCubic'
    });

   
     
   
    $(".hamburger-menu, .main-nav ul li a").on( 'click', function() {
        $(".header").toggleClass("pushed");
        $(".main-content").toggleClass("main-pushed");
        $('.bar').toggleClass('animate');
    });

   
   

   
      $(".customer-carousel").owlCarousel({
        items: 4
      });

      $(".resume-carousel, .testimonial-carousel").owlCarousel({
        singleItem:true
      });

   

   
    
    $(".header").css({'height':($(window).height())+'px'});

    $(document).ready(function () {

      $(window).scroll(function () {
        if ($(this).scrollTop() > 100) {
          $('.scrollup').fadeIn();
        } else {
          $('.scrollup').fadeOut();
        }
      });

      $('.scrollup').click(function () {
        $("html, body, main").animate({
          scrollTop: 0
        }, 600);
        return false;
      });

    });