$(document).ready(function(){
    $(".side-bar-click").click(function(){
      $(".side-menu").addClass("sidebar-active");
    });

    $(document).click(function(){
        $('.side-menu').removeClass('sidebar-active');
      });   
      $("iframe , .white-cross").click(function(){
        $('.side-menu').removeClass('sidebar-active');
      });
      $(".side-bar-click, .side-menu  ").click(function(e){
        e.stopPropagation();
      });
  });

  $(document).ready(function(){
    $(".map-start-circle .btn").click(function(){
      $(".rate-card-popu").addClass("rate-card-popu-active");
    });
    $(document).click(function(){
      $('.rate-card-popu ').removeClass('rate-card-popu-active');
    });   
    $("iframe").click(function(){
      $('.rate-card-popu').removeClass('rate-card-popu-active');
    });
    $(".map-start-circle .btn  ,.rate-card-popu ,#myModal").click(function(e){
      e.stopPropagation();
    });
  });

  $(document).ready(function(){
  $(".sel-dropdown").change(function(){
    $(this).attr('data-option','hidden')
    if($(this).val()){
        $(this).removeAttr('data-option')
    }
})
});



